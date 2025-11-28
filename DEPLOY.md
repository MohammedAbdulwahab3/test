# Deploy to Render and Grafana Cloud

Complete guide for deploying the Flutter Go CRUD backend to Render (free tier) with Grafana Cloud monitoring (free tier).

## Prerequisites

1. **Render Account** (Free): Sign up at https://render.com
2. **Grafana Cloud Account** (Free): Sign up at https://grafana.com/auth/sign-up/create-user
3. **GitHub Account**: Code must be in a GitHub repository

## Part 1: Deploy to Render

### Step 1: Push Code to GitHub

```bash
cd /home/maw/Desktop/try/flutter_go_crud/backend

# Initialize git if not already done
git init

# Add all files
git add .

# Commit
git commit -m "Add Prometheus metrics and Render deployment"

# Create GitHub repo and push
git remote add origin https://github.com/YOUR_USERNAME/flutter-backend.git
git branch -M main
git push -u origin main
```

### Step 2: Create Render Account

1. Go to https://render.com
2. Sign up with GitHub account (recommended)
3. Authorize Render to access your GitHub repositories

### Step 3: Deploy Using Blueprint

1. **Go to Render Dashboard**: https://dashboard.render.com
2. **Click "New" → "Blueprint"**
3. **Connect your GitHub repository**
4. **Select the repository** containing your backend code
5. **Render will automatically detect** `render.yaml`
6. **Review the services**:
   - `flutter-backend-db`: PostgreSQL database (Free, 90 days)
   - `flutter-backend`: Web service (Free tier)
7. **Click "Apply"**

### Step 4: Wait for Deployment

- Database provisioning: ~2-3 minutes
- Backend build & deploy: ~5-10 minutes
- Watch logs in real-time on Render dashboard

### Step 5: Get Your Backend URL

Once deployed, you'll get a URL like:
```
https://flutter-backend-XXXX.onrender.com
```

### Step 6: Test Deployment

```bash
# Health check
curl https://YOUR-APP.onrender.com/healthz

# Metrics endpoint
curl https://YOUR-APP.onrender.com/metrics

# API test
curl https://YOUR-APP.onrender.com/api/items
```

## Part 2: Setup Grafana Cloud Monitoring

### Step 1: Create Grafana Cloud Account

1. Go to https://grafana.com/auth/sign-up/create-user
2. Fill in your details
3. Verify email
4. Choose **Free tier** (10k metrics, 50GB logs, 14-day retention)

### Step 2: Get Prometheus Remote Write URL

1. In Grafana Cloud dashboard, go to **"My Account"**
2. Click on your stack
3. Go to **"Details & API Keys"**
4. Find **"Prometheus"** section
5. Note the **Remote Write Endpoint**: `https://prometheus-XXXX.grafana.net/api/prom/push`
6. Generate an API key if you don't have one

### Step 3: Configure Prometheus Agent (Simple Option)

Since we're using Render free tier, the easiest approach is to use **Grafana Agent** or **Grafana Alloy**:

1. **Add Grafana Alloy to render.yaml**:

```yaml
# Add this to render.yaml
- type: web
  name: grafana-alloy
  env: docker
  plan: free
  dockerfilePath: ./Dockerfile.alloy
  envVars:
    - key: PROMETHEUS_URL
      value: https://prometheus-XXXX.grafana.net/api/prom/push
    - key: PROMETHEUS_USER
      value: YOUR_INSTANCE_ID
    - key: PROMETHEUS_API_KEY
      value: YOUR_API_KEY
    - key: SCRAPE_TARGET
      value: https://flutter-backend-XXXX.onrender.com/metrics
```

2. **Create Dockerfile.alloy**:

```dockerfile
FROM grafana/alloy:latest
COPY alloy-config.yaml /etc/alloy/config.yaml
CMD ["run", "--server.http.listen-addr=0.0.0.0:12345", "/etc/alloy/config.yaml"]
```

3. **Create alloy-config.yaml**:

```yaml
prometheus.scrape "backend" {
  targets = [{
    __address__ = env("SCRAPE_TARGET"),
  }]
  forward_to = [prometheus.remote_write.grafana_cloud.receiver]
  scrape_interval = "15s"
}

prometheus.remote_write "grafana_cloud" {
  endpoint {
    url = env("PROMETHEUS_URL")
    basic_auth {
      username = env("PROMETHEUS_USER")
      password = env("PROMETHEUS_API_KEY")
    }
  }
}
```

### Alternative: Manual Metrics Forwarding

If you want a simpler setup without Grafana Alloy:

1. Use **Grafana Cloud's HTTP API** to push metrics
2. Or use the `/metrics` endpoint directly from Grafana (requires authentication setup)

### Step 4: Import Dashboard

1. In Grafana Cloud, go to **"Dashboards"**
2. Click **"Import"**
3. Click **"Upload JSON file"**
4. Upload `grafana/dashboards/backend-overview.json`
5. Select your Prometheus data source
6. Click **"Import"**

### Step 5: Verify Metrics

1. Go to your imported dashboard
2. You should see:
   - HTTP Request Rate
   - Request Latency (p50, p95, p99)
   - Error Rates (4xx, 5xx)
   - Items by Status
   - CRUD Operations Rate
   - Database Query Duration

## Part 3: Testing Everything

### Generate Some Traffic

```bash
# Create items
for i in {1..10}; do
  curl -X POST https://YOUR-APP.onrender.com/api/items \
    -H "Content-Type: application/json" \
    -d "{\"name\":\"Item $i\",\"description\":\"Test item\",\"price\":99.99,\"status\":\"active\"}"
done

# List items
curl https://YOUR-APP.onrender.com/api/items

# Check metrics
curl https://YOUR-APP.onrender.com/metrics | grep http_requests_total
```

### View in Grafana

1. Go to your Grafana dashboard
2. Set time range to "Last 5 minutes"
3. You should see metrics appearing:
   - Request rate increasing
   - Latency graphs
   - CRUD operations counter

## Important Notes

### Render Free Tier Limitations

- 🌙 **Sleeps after 15 minutes of inactivity**
- ⚡ **First request after sleep takes ~30-60 seconds** (cold start)
- 🔄 **750 hours/month** (auto-sleep helps stay within limit)
- 💾 **PostgreSQL free for 90 days**, then $7/month

### Grafana Cloud Free Tier

- 📊 **10,000 series** (metrics)
- 📝 **50GB logs**
- ⏱️ **14-day retention**
- ✅ **Perfect for hobby projects**

### Keeping Service Awake (Optional)

If you want to avoid cold starts:

1. Use a service like **UptimeRobot** (free)
2. Ping your `/healthz` endpoint every 5 minutes
3. This keeps the service active

**Note**: This will consume your 750 hours faster, use wisely!

## Costs Summary

| Service | Free Tier | After Free |
|---------|-----------|------------|
| **Render Web Service** | 750 hours/month | $7/month (Starter) |
| **Render PostgreSQL** | 90 days | $7/month |
| **Grafana Cloud** | Forever free | Paid plans available |

**Total for testing**: **$0/month** ✅

## Troubleshooting

### Build Fails

- Check Render logs
- Ensure `go.mod` and `go.sum` are committed
- Verify `Dockerfile` is correct

### Metrics Not Showing

- Verify `/metrics` endpoint returns data
- Check Grafana Alloy logs
- Ensure API key is correct
- Check Prometheus remote write URL

### Database Connection Issues

- Check environment variables in Render
- Verify `DB_DSN` is set from database connection
- Check database logs

### Cold Starts Too Slow

- Use UptimeRobot to ping every 5 min
- Or upgrade to paid tier ($7/month, no sleep)

## Next Steps

1. ✅ Deploy backend to Render
2. ✅ Setup Grafana Cloud monitoring
3. 📱 Update Flutter frontend with Render URL
4. 🎨 Customize Grafana dashboard
5. 🔔 Setup alerts (optional, available in Grafana Cloud free)
6. 📊 Add more custom metrics as needed

## Useful Links

- **Render Dashboard**: https://dashboard.render.com
- **Grafana Cloud**: https://grafana.com
- **Your Backend**: https://YOUR-APP.onrender.com
- **Metrics**: https://YOUR-APP.onrender.com/metrics
- **API Docs**: https://YOUR-APP.onrender.com/api/items
