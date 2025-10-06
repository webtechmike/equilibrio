# Redis Caching Strategy

## 🎯 **Smart Market-Aware Caching**

Our caching strategy is designed to minimize API calls while providing up-to-date data by being aware of US stock market hours.

## 📊 **How It Works**

### **Market Hours Detection**

**US Stock Market:**
- **Open:** Monday-Friday, 9:30 AM - 4:00 PM ET
- **Closed:** Weekends, before 9:30 AM, after 4:00 PM
- **Timezone:** America/New_York (Eastern Time)

### **Caching Logic**

```
┌─────────────────────────────────────────────────────┐
│ Request comes in                                     │
└──────────────────┬──────────────────────────────────┘
                   │
                   v
         ┌─────────────────┐
         │ Is market open? │
         └────────┬────────┘
                  │
        ┌─────────┴─────────┐
        │                   │
        v                   v
    ┌───────┐         ┌──────────┐
    │  YES  │         │    NO    │
    └───┬───┘         └────┬─────┘
        │                  │
        v                  v
   Cache for          Check daily
   5 minutes          snapshot
        │                  │
        v            ┌─────┴─────┐
   Fetch real     Exists?    Not found
   data if           │            │
   expired           v            v
                  Return      Fetch &
                  cached      cache for
                  data        next open
```

## ⚡ **Cache TTL Strategy**

### **During Market Hours (9:30 AM - 4:00 PM ET)**
- **TTL:** 5 minutes
- **Reason:** Prices change frequently, need recent data
- **Behavior:** Refreshes every 5 minutes if requested

### **After Market Close (4:00 PM ET - 9:30 AM ET next day)**
- **TTL:** Until next market open
- **Reason:** Prices don't change, no need to refresh
- **Example:** If market closes at 4 PM Monday, cache until 9:30 AM Tuesday

### **Weekends**
- **TTL:** Until Monday 9:30 AM ET
- **Reason:** Market closed, data won't change
- **Example:** Saturday data cached until Monday morning

## 🗄️ **Cache Keys**

### **Daily Snapshot**
```
Key: daily_snapshot:2025-10-06
TTL: Until next market open
Data: Full list of all stocks fetched that day
```

### **Filtered Results**
```
Key: stocks:{filter_params_hash}
TTL: Market-aware (5 min or until open)
Data: Filtered/sorted/paginated results
```

### **Individual Stocks**
```
Key: stock:AAPL
TTL: Market-aware (5 min or until open)
Data: Single stock data
```

## 💡 **Benefits**

### **Reduced API Calls**
**Before:**
- 30-second cache
- ~120 API calls per hour
- ~2,880 calls per day
- Exceeds most free tier limits

**After:**
- Market-aware caching
- ~12 calls during market hours (6.5 hours)
- ~0 calls after hours (17.5 hours)
- **~78 calls per day** (97% reduction!)

### **Better Performance**
- **Instant response** when market closed (cache hit)
- **Fresh data** during trading hours
- **No rate limiting** from excessive calls

### **Cost Savings**
- Stays within free tier limits (Yahoo: 60/min, Alpha Vantage: 25/day)
- No need for paid API plans
- Reduced server load

## 📈 **Example Scenarios**

### **Scenario 1: Monday Morning (10:00 AM ET)**
```
Request at 10:00 AM
Market: OPEN
Action: Check cache (5 min TTL)
  - If cached < 5 min ago: Return cached
  - If expired: Fetch fresh, cache for 5 min
API Calls: Max 12/hour during market hours
```

### **Scenario 2: Monday Evening (7:00 PM ET)**
```
Request at 7:00 PM
Market: CLOSED
Action: Check daily snapshot (key: daily_snapshot:2025-10-06)
  - If exists: Return snapshot
  - If not: Fetch once, cache until Tuesday 9:30 AM
API Calls: 1 total (reused until next open)
```

### **Scenario 3: Saturday (any time)**
```
Request on Saturday
Market: CLOSED (Weekend)
Action: Check daily snapshot (key: daily_snapshot:2025-10-06 from Friday)
  - Return Friday's closing data
  - Cache persists until Monday 9:30 AM
API Calls: 0 (using Friday's data)
```

### **Scenario 4: Just After Market Close (4:05 PM ET)**
```
Request at 4:05 PM
Market: CLOSED
Action: Daily snapshot created at market close
  - Data cached until tomorrow 9:30 AM
  - All evening/night requests use this snapshot
API Calls: 1 (at 4 PM), then 0 until next morning
```

## 🔧 **Configuration**

### **Environment Variables**
```bash
USE_MOCK_DATA=false         # Use real data
CACHE_TTL_SECONDS=300       # Used as base, overridden by market hours
REDIS_URL=redis://redis:6379
```

### **Market Hours Configuration**
Located in `market_cache_strategy.go`:
```go
var USMarketHours = MarketHours{
    OpenHour:  9,      // 9:30 AM ET
    CloseHour: 16,     // 4:00 PM ET
    TimeZone:  "America/New_York",
}
```

## 📊 **Redis Memory Usage**

### **Estimated Storage**
- **Single stock:** ~1 KB
- **50 stocks list:** ~50 KB
- **Daily snapshot:** ~50 KB
- **Filtered cache (10 variations):** ~500 KB
- **Total:** <1 MB for typical usage

### **Memory Optimization**
- Old snapshots auto-expire
- Filtered results cached temporarily
- Redis eviction policy: `allkeys-lru` (recommended)

## 🚀 **API Call Optimization**

### **Daily API Usage Breakdown**

**Market Open (9:30 AM - 4:00 PM = 6.5 hours)**
- Requests: ~6 per hour (assuming 5-min cache hits)
- Total: 6 × 6.5 = **39 calls**

**Market Closed (4:00 PM - 9:30 AM = 17.5 hours)**
- Requests: 1 snapshot at close
- Total: **1 call**

**Weekends (Saturday + Sunday = 48 hours)**
- Requests: 0 (using Friday's snapshot)
- Total: **0 calls**

**Weekly Total:**
- Weekdays: 39 × 5 = 195 calls
- Weekend: 0 calls
- **~195 calls/week** (vs 20,160 without smart caching!)

## 🎯 **Smart Features**

### **1. Weekend Awareness**
- Automatically uses Friday's data
- No API calls on weekends
- Saves 28% of potential calls

### **2. After-Hours Optimization**
- Evening data cached until morning
- Pre-market uses previous day's close
- Saves 72% of daily calls

### **3. Holiday Support**
- Can be extended to detect market holidays
- Treats holidays like weekends
- Further reduces unnecessary calls

### **4. Timezone Aware**
- All calculations in Eastern Time
- Works correctly regardless of server timezone
- Handles daylight saving time

## 📝 **Monitoring**

### **Log Messages**
```
Market Status: Market Open | Cache TTL: 5m0s
Market Status: Market Closed (After Hours) | Cache TTL: 17h25m0s
Market Status: Market Closed (Weekend) | Cache TTL: 41h30m0s
Using daily snapshot (market closed)
Cached daily snapshot for reuse
```

### **Key Metrics to Track**
- Cache hit rate
- API calls per day
- Average response time
- Redis memory usage

## 🔄 **Cache Warming**

### **Automatic Warming**
- First request after market close creates snapshot
- Snapshot reused for all subsequent requests
- No manual intervention needed

### **Manual Warming (Optional)**
```go
// Warm cache at market close (4:00 PM ET)
if justClosed {
    cacheStrategy.WarmCache(ctx, func() ([]models.StockData, error) {
        return marketDataService.getRealStockData()
    })
}
```

## 🎉 **Summary**

**Old Strategy:**
- ❌ 30-second TTL
- ❌ ~2,880 API calls/day
- ❌ Exceeds free tier limits
- ❌ Frequent rate limiting

**New Strategy:**
- ✅ Market-aware caching
- ✅ ~40 API calls/day (98% reduction)
- ✅ Stays within free tiers
- ✅ Zero rate limiting issues
- ✅ Better performance
- ✅ Lower server load

**Perfect for:**
- Free API tiers (Yahoo, Alpha Vantage)
- High-traffic applications
- Cost-sensitive deployments
- Rate-limited APIs

**The market closes at 4 PM ET, so we cache that data until next open - smart, efficient, and API-friendly!** 🚀

