import { ref } from 'vue'
import { api } from '../services/api'

export function useAIAnalysis() {
  const analysis = ref(null)
  const isAnalyzing = ref(false)
  const error = ref(null)

  const analyzeToken = async (tokenData) => {
    isAnalyzing.value = true
    error.value = null
    analysis.value = null
    
    try {
      // Prepare minimal data for AI to reduce token count
      const payload = {
         symbol: tokenData.symbol,
         name: tokenData.name,
         price: tokenData.price,
         market_cap: tokenData.market_cap,
         volume_24h: tokenData.volume_24h,
         change_24h: tokenData.change_24h,
         change_7d: tokenData.change_7d,
         change_30d: tokenData.change_30d || 0,
         change_90d: tokenData.change_90d || 0,
         tvl: tokenData.tvl || 0,
         trust_score: tokenData.trust_score,
         liquidity: tokenData.liquidity || 0,
         rank: tokenData.rank || 0,
         holder_count: tokenData.holder_count || 0,
         circulating_supply: tokenData.circulating_supply || 0,
         max_supply: tokenData.max_supply || 0,
         total_supply: tokenData.total_supply || 0
      }
      
      const result = await api.analyzeToken(payload)
      
      if (result && result.analysis) {
         analysis.value = result.analysis
      } else if (result && result.status === 'success' && result.analysis) {
         analysis.value = result.analysis
      } else {
         console.warn('AI response missing analysis field:', result)
         error.value = "AI response was empty. Please check backend logs."
      }
    } catch (e) {
      console.error('AI Analysis failed:', e)
      error.value = "Failed to generate analysis. Please try again."
    } finally {
      isAnalyzing.value = false
    }
  }

  const clearAnalysis = () => {
    analysis.value = null
    error.value = null
  }

  return {
    analysis,
    isAnalyzing,
    error,
    analyzeToken,
    clearAnalysis
  }
}
