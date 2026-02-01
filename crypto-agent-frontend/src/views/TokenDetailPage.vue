<template>
  <div class="min-h-screen bg-[#050510]">
    <!-- Full Page Container -->
    <div class="p-0 overflow-hidden bg-[#0a0b0f] w-full mx-auto flex flex-col min-h-screen relative text-gray-300">
      
      <!-- HERO HEADER SECTION -->
      <div class="relative p-10 pb-16 overflow-hidden">
        <!-- Dynamic Hero Background -->
        <div class="absolute inset-0 bg-gradient-to-br from-[#1a1c24] via-[#0a0b10] to-[#050608] z-0"></div>
        <div class="absolute top-0 right-0 w-[600px] h-[600px] bg-primary/10 blur-[120px] rounded-full translate-x-1/2 -translate-y-1/2 z-0 animate-pulse"></div>
        
        <!-- Back Button -->
        <button @click="$router.push('/')" class="absolute top-8 left-8 p-3 rounded-2xl bg-white/5 hover:bg-white/10 border border-white/10 text-white transition-all z-20 hover:scale-110 active:scale-95 group flex items-center gap-2">
          <svg class="w-6 h-6 group-hover:-translate-x-1 transition-transform duration-300" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M15 19l-7-7 7-7"/></svg>
          <span class="text-sm font-bold uppercase tracking-wider">Back</span>
        </button>

        <div v-if="displayToken && displayToken.symbol" class="relative z-10">
          <div class="flex flex-col lg:flex-row lg:items-center justify-between gap-10">
            <!-- Left: Identity -->
            <div class="flex items-center gap-8">
              <div class="relative group">
                <div class="absolute -inset-4 bg-primary/20 blur-2xl rounded-full opacity-0 group-hover:opacity-100 transition-opacity"></div>
                <div class="w-24 h-24 rounded-3xl bg-gradient-to-br from-gray-800 to-black flex items-center justify-center border border-white/20 shadow-2xl relative z-10 overflow-hidden backdrop-blur-xl">
                  <img v-if="displayToken.image" :src="displayToken.image" :alt="displayToken.name" class="w-full h-full object-cover scale-90 group-hover:scale-100 transition-transform duration-500" />
                  <span v-else class="text-4xl font-black text-primary">{{ displayToken.symbol?.[0] || '?' }}</span>
                </div>
              </div>
              
              <div class="space-y-3">
                <div class="flex items-center gap-4">
                  <h2 class="text-6xl font-black text-white tracking-tighter leading-tight">{{ displayToken.name }}</h2>
                  <span class="px-3 py-1 bg-white/5 border border-white/10 rounded-lg text-gray-400 text-sm font-mono font-bold tracking-widest uppercase h-min mt-2">{{ displayToken.symbol }}</span>
                </div>
                <div class="flex items-center gap-3">
                  <div class="flex items-center gap-2 px-4 py-1.5 bg-primary/10 text-primary text-xs font-black rounded-full border border-primary/20 shadow-lg shadow-primary/5">
                    <span class="w-1.5 h-1.5 bg-primary rounded-full animate-pulse"></span>
                    RANK #{{ displayToken.rank }}
                  </div>
                  <div v-if="displayToken.category" class="px-4 py-1.5 bg-white/5 text-gray-400 text-xs font-black rounded-full border border-white/10 uppercase tracking-widest">
                    {{ displayToken.category }}
                  </div>
                </div>
              </div>
            </div>

            <!-- Right: Real-time Wealth Panel -->
            <div class="flex flex-col items-end gap-2 text-right">
              <div class="text-[10px] font-black text-primary uppercase tracking-[0.4em] mb-1">Live Asset Valuation</div>
              <div class="text-7xl font-black text-white tracking-tighter mb-2 tabular-nums drop-shadow-2xl">
                {{ formatCurrency(displayToken.price) }}
              </div>
              
              <div class="flex items-center gap-4 justify-end">
                <!-- 24H Change Badge -->
                <div 
                  class="flex items-center gap-2 px-4 py-2 rounded-2xl font-black text-sm border shadow-lg"
                  :class="displayToken.change_24h >= 0 ? 'bg-green-500/10 border-green-500/30 text-green-400' : 'bg-red-500/10 border-red-500/30 text-red-500'"
                >
                  <span class="text-xl leading-none">{{ displayToken.change_24h >= 0 ? '↗' : '↘' }}</span>
                  <span class="font-mono">{{ Math.abs(displayToken.change_24h).toFixed(2) }}%</span>
                  <span class="text-[10px] opacity-40 uppercase ml-1">24H</span>
                </div>
                
                <!-- 7D Change Badge -->
                <div 
                  class="flex items-center gap-2 px-4 py-2 rounded-2xl font-black text-sm border shadow-lg"
                  :class="displayToken.change_7d >= 0 ? 'bg-green-500/10 border-green-500/20 text-green-400/80' : 'bg-red-500/10 border-red-500/20 text-red-500/80'"
                >
                   <span class="font-mono">{{ Math.abs(displayToken.change_7d).toFixed(2) }}%</span>
                   <span class="text-[10px] opacity-40 uppercase ml-1">7D</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Quick Stats Belt -->
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mt-12 pt-8 border-t border-white/5 relative z-10">
            <div class="group px-2">
              <div class="text-[10px] text-gray-500 font-bold uppercase tracking-[0.2em] mb-2 group-hover:text-primary transition-colors">Circulating Supply</div>
              <div class="text-lg font-mono text-white flex items-end gap-2">
                {{ formatNumber(displayToken.circulating_supply) }} 
                <span class="text-[10px] text-gray-600 mb-1 font-sans uppercase">{{ displayToken.symbol }}</span>
              </div>
              <div class="w-full h-1 bg-white/5 rounded-full mt-2 overflow-hidden">
                <div class="h-full bg-primary/40" :style="{ width: displayToken.max_supply > 0 ? (displayToken.circulating_supply / displayToken.max_supply * 100) + '%' : '100%' }"></div>
              </div>
            </div>
            <div class="group px-2">
              <div class="text-[10px] text-gray-500 font-bold uppercase tracking-[0.2em] mb-2 group-hover:text-amber-500 transition-colors">Max Token Supply</div>
              <div class="text-lg font-mono text-white">{{ displayToken.max_supply > 0 ? formatNumber(displayToken.max_supply) : '∞ UNLIMITED' }}</div>
            </div>
            <div class="group px-2">
              <div class="text-[10px] text-gray-500 font-bold uppercase tracking-[0.2em] mb-2 group-hover:text-cyan-500 transition-colors">Market Share</div>
              <div class="text-lg font-mono text-cyan-400 font-bold tracking-widest">{{ (displayToken.MarketCapDom || displayToken.market_cap_dominance || 0).toFixed(3) }}%</div>
            </div>
            <div class="group px-2">
              <div class="text-[10px] text-gray-500 font-bold uppercase tracking-[0.2em] mb-2 group-hover:text-green-500 transition-colors">Peak Valuation (FDV)</div>
              <div class="text-lg font-mono text-white tracking-tight">{{ formatCurrencyCompact(displayToken.fdv) }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- MAIN APP BODY -->
      <div class="flex-1 overflow-y-auto custom-scrollbar bg-[#050608] p-10 rounded-t-[60px] border-t border-white/10 relative z-10 -mt-10 shadow-[0_-50px_100px_rgba(0,0,0,0.5)]">
        <div class="grid lg:grid-cols-12 gap-10">
          
          <!-- LEFT COLUMN: Advanced Scoring & Metrics (5 cols) -->
          <div class="lg:col-span-4 space-y-8">
            <!-- Global Alpha Score Card -->
            <div class="glass-bento p-10 flex flex-col items-center group relative overflow-hidden">
              <div class="absolute -top-10 -left-10 w-40 h-40 bg-primary/10 blur-[60px] rounded-full group-hover:bg-primary/20 transition-all duration-700"></div>
              <div class="text-[10px] font-black text-gray-500 uppercase tracking-[0.4em] mb-8 w-full text-center">Core Trust Rating</div>
              
              <div class="relative mb-8">
                <div class="absolute inset-0 bg-primary/20 blur-3xl rounded-full scale-150 opacity-40"></div>
                <RatingBadge 
                  :grade="displayToken.score_breakdown?.grade || 'D'" 
                  class="text-6xl px-10 py-6 rounded-[32px] shadow-[0_20px_60px_rgba(0,0,0,0.5)] border-2 border-white/10 relative z-10"
                />
              </div>
              
              <div class="text-center group-hover:scale-105 transition-transform">
                <div class="text-7xl font-black text-white tracking-tighter leading-none">{{ displayToken.trust_score?.toFixed(0) || '—' }}</div>
                <div class="text-xs text-primary font-black uppercase tracking-[0.3em] mt-3">Alpha Index</div>
              </div>
            </div>

            <!-- Sentiment & Vibe -->
            <div class="glass-bento p-10 text-center">
              <div class="text-[10px] font-black text-gray-500 uppercase tracking-[0.4em] mb-10 text-left">Community Vibe</div>
              <div class="flex flex-col items-center justify-center py-4">
                <SentimentGauge :score="displayToken.sentiment || 0" class="scale-125" />
                <div class="mt-12 text-[11px] font-black text-gray-400 bg-white/5 px-6 py-3 rounded-2xl border border-white/5 flex items-center gap-3">
                  <span class="w-2 h-2 rounded-full" :class="displayToken.sentiment >= 0 ? 'bg-green-400 shine' : 'bg-rose-400 shine'"></span>
                  SENTIMENT SHIFT: <span :class="displayToken.sentiment >= 0 ? 'text-green-400' : 'text-rose-400'">{{ (displayToken.sentiment * 100).toFixed(1) }}% MATCH</span>
                </div>
              </div>
            </div>

            <!-- Radar Matrix (Scaled Up) -->
            <div class="glass-bento p-10 flex flex-col items-center justify-center min-h-[480px]">
              <div class="text-[10px] font-black text-gray-500 uppercase tracking-[0.4em] mb-10 w-full text-left">Asset Signature Matrix</div>
              <div class="w-full h-[380px] relative">
                <ScoreRadar v-if="displayToken.score_breakdown" :scores="{
                  liquidity: displayToken.score_breakdown.liquidity_score || 0,
                  volume: displayToken.score_breakdown.volume_score || 0,
                  trend: displayToken.score_breakdown.trend_score || 0,
                  social: displayToken.score_breakdown.social_score || 0,
                  market: displayToken.score_breakdown.market_health_score || 0,
                  stability: displayToken.score_breakdown.risk_score || 0
                }" />
                <div v-else class="flex items-center justify-center h-full text-gray-600 font-black uppercase tracking-[0.2em] italic text-sm animate-pulse">
                  Calibrating Matrix...
                </div>
              </div>
            </div>
          </div>

          <!-- RIGHT COLUMN: Data Grid & AI Analysis (7 cols) -->
          <div class="lg:col-span-8 space-y-10">
            
            <!-- Section: Institutional Benchmarks -->
            <div class="space-y-4">
              <h4 class="text-[10px] font-black text-primary uppercase tracking-[0.4em] px-4">Financial Benchmarks</h4>
              <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                <MetricCard 
                  label="Capitalization" 
                  :value="formatCurrencyCompact(displayToken.market_cap)" 
                  icon="💎"
                  :subLabel="`GLOBAL RANK #${displayToken.rank}`"
                  class="bg-gradient-to-br from-white/[0.03] to-transparent hover:border-primary/30 transition-all duration-500"
                />
                <MetricCard 
                   label="Velocity (24V)" 
                   :value="formatCurrencyCompact(displayToken.volume_24h)" 
                   icon="⚡"
                   :subLabel="displayToken.volume_change_24h ? (displayToken.volume_change_24h >= 0 ? '↑ ' : '↓ ') + displayToken.volume_change_24h.toFixed(1) + '% Momentum' : 'N/A Momentum'"
                   :progress="displayToken.market_cap > 0 ? Math.min((displayToken.volume_24h / displayToken.market_cap) * 200, 100) : 0"
                />
                <MetricCard 
                  label="Peak Retracement" 
                  :value="displayToken.ath_change ? displayToken.ath_change.toFixed(1) + '%' : '—'" 
                  :change="displayToken.ath_change"
                  icon="📉"
                  :progress="Math.abs(displayToken.ath_change || 0)"
                  progressColor="rose"
                />
              </div>
            </div>

            <!-- Section: Advanced On-Chain Health -->
            <div class="space-y-4">
              <h4 class="text-[10px] font-black text-amber-500/80 uppercase tracking-[0.4em] px-4">Protocol Vital Signs</h4>
              <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                <MetricCard 
                   label="Network Holders" 
                   :value="displayToken.holder_count > 0 ? formatNumber(displayToken.holder_count) : 'SCANNED'" 
                   icon="👥"
                   subLabel="Verified Addresses"
                   class="hover:border-amber-500/30 transition-all duration-500"
                />
                <MetricCard 
                   label="Depth / Liquidity" 
                   :value="displayToken.liquidity > 0 ? formatCurrencyCompact(displayToken.liquidity) : 'LIMITED'" 
                   icon="💧"
                   :subLabel="displayToken.market_cap > 0 ? ((displayToken.liquidity / displayToken.market_cap) * 100).toFixed(2) + '% Ratio' : 'Liq Density'"
                   :progress="displayToken.liquidity > 0 ? Math.min((displayToken.liquidity / displayToken.market_cap) * 1000, 100) : 0"
                />
                <MetricCard 
                   label="Network Revenue" 
                   :value="displayToken.revenue_30d > 0 ? formatCurrencyCompact(displayToken.revenue_30d) : 'EARLY PHASE'" 
                   icon="💰"
                   :subLabel="displayToken.ps_ratio > 0 ? 'P/S SCALE: ' + displayToken.ps_ratio.toFixed(1) : 'Fee Protocol'"
                   progressColor="emerald"
                />
              </div>
            </div>

            <!-- AI Alpha Agent Insight -->
            <div class="relative group mt-4">
              <div class="absolute -inset-1 bg-gradient-to-r from-primary/30 via-indigo-500/20 to-primary/30 rounded-[40px] blur-2xl opacity-50 group-hover:opacity-100 transition duration-1000"></div>
              
              <div class="glass-bento p-12 relative overflow-hidden bg-black/80 backdrop-blur-3xl border border-white/10 rounded-[40px] shadow-2xl">
                <!-- Decorative Elements -->
                <div class="absolute top-0 right-0 w-full h-full bg-[radial-gradient(circle_at_100%_0%,rgba(59,130,246,0.1)_0%,transparent_50%)]"></div>
                <div class="absolute -bottom-24 -left-24 w-96 h-96 bg-primary/5 blur-[120px] rounded-full"></div>
                
                <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-6 mb-12 relative z-10">
                   <div class="flex items-center gap-5">
                     <div class="w-16 h-16 rounded-[22px] bg-primary/10 flex items-center justify-center text-primary text-3xl shadow-lg border border-primary/20 group-hover:scale-110 transition-transform duration-500">✨</div>
                     <div>
                       <h3 class="text-2xl font-black text-white tracking-widest uppercase">Agent Alpha Report</h3>
                       <div class="text-[10px] text-primary font-black uppercase tracking-[0.4em] flex items-center gap-2 mt-2">
                          <span class="w-2 h-2 bg-primary rounded-full animate-ping"></span>
                          Full Spectrum Synthesis Active
                       </div>
                     </div>
                   </div>
                   <div v-if="isAnalyzing" class="flex items-center gap-4 px-6 py-3 bg-primary/10 border border-primary/20 rounded-2xl backdrop-blur-md">
                     <div class="flex gap-1.5">
                        <span class="w-2 h-4 bg-primary rounded-full animate-bounce [animation-delay:-0.3s]"></span>
                        <span class="w-2 h-4 bg-primary rounded-full animate-bounce [animation-delay:-0.15s]"></span>
                        <span class="w-2 h-4 bg-primary rounded-full animate-bounce"></span>
                     </div>
                     <span class="text-xs font-black text-primary uppercase tracking-[0.2em] animate-pulse">Processing Market Signals...</span>
                   </div>
                </div>
                
                <!-- Structured JSON AI UI -->
                <div v-if="parsedAnalysis" class="space-y-16 relative z-10">
                  <!-- Executive Summary -->
                  <div class="relative">
                    <div class="absolute -left-8 top-0 bottom-0 w-1.5 bg-gradient-to-b from-primary to-transparent rounded-full opacity-60"></div>
                    <p class="text-3xl font-bold text-white leading-tight tracking-tight pl-6 pr-4">
                      "{{ parsedAnalysis.summary }}"
                    </p>
                  </div>

                  <div class="grid md:grid-cols-2 gap-10">
                    <!-- Potential Score UI -->
                    <div class="p-10 rounded-[36px] bg-white/[0.03] border border-white/5 space-y-6 group/potential hover:bg-white/[0.05] transition-all">
                      <div class="flex justify-between items-end">
                         <span class="text-xs font-black text-gray-400 uppercase tracking-[0.4em]">Future Growth Potential</span>
                         <span class="text-5xl font-black text-primary tabular-nums drop-shadow-[0_0_15px_rgba(59,130,246,0.5)]">
                           {{ parsedAnalysis.growth_potential.score }}<span class="text-lg opacity-50 ml-1 font-bold">%</span>
                         </span>
                      </div>
                      <div class="h-4 w-full bg-black/60 rounded-full overflow-hidden border border-white/10 p-0.5">
                        <div class="h-full rounded-full bg-gradient-to-r from-primary via-indigo-400 to-primary bg-[length:200%_100%] animate-shimmer transition-all duration-1500" :style="{ width: `${parsedAnalysis.growth_potential.score}%` }"></div>
                      </div>
                      <p class="text-sm text-gray-400 leading-relaxed font-medium mt-4">{{ parsedAnalysis.growth_potential.reason }}</p>
                    </div>

                    <!-- Recommendation Card -->
                    <div class="p-10 rounded-[36px] bg-gradient-to-br from-primary/20 via-primary/5 to-transparent border border-primary/20 shadow-xl transition-all group/recommend hover:border-primary/40">
                      <div class="text-[10px] font-black text-primary uppercase tracking-[0.4em] mb-6">Strategic Recommendation</div>
                      <div class="flex items-center justify-between mb-8">
                         <span class="text-5xl font-black text-white tracking-widest uppercase">{{ parsedAnalysis.recommendation.action }}</span>
                         <div class="text-[11px] font-black text-primary bg-primary/20 px-4 py-2 rounded-xl border border-primary/30 tracking-[0.2em] font-mono shadow-lg">ZONE: {{ parsedAnalysis.recommendation.entry_zone }}</div>
                      </div>
                      <div class="pt-6 border-t border-white/10 text-xs text-gray-500 font-bold uppercase tracking-[0.3em]">Target Exit: <span class="text-white ml-2 text-sm">{{ parsedAnalysis.recommendation.target }}</span></div>
                    </div>
                  </div>

                  <!-- Technical & Risk row -->
                  <div class="grid md:grid-cols-3 gap-6">
                    <div class="p-8 rounded-3xl bg-white/[0.03] border border-white/5 flex flex-col justify-between hover:border-primary/20 transition-all">
                      <div class="text-[10px] font-black text-gray-500 uppercase tracking-[0.4em] mb-4">Technical Bias</div>
                      <div class="text-2xl font-black text-white uppercase tracking-tight">{{ parsedAnalysis.technical_analysis.trend }}</div>
                      <div class="text-[11px] text-primary font-black uppercase tracking-[0.3em] mt-3 flex items-center gap-2">
                        <span class="w-1.5 h-1.5 rounded-full bg-primary shadow-glow"></span>
                        {{ parsedAnalysis.technical_analysis.strength }} STRENGTH
                      </div>
                    </div>
                    
                    <div class="p-8 rounded-3xl border flex flex-col justify-between hover:scale-[1.02] transition-all" :class="riskClass">
                      <div class="text-[10px] font-black uppercase tracking-[0.4em] mb-4 opacity-70">Risk Assessment</div>
                      <div class="text-2xl font-black uppercase tracking-tight">{{ parsedAnalysis.risk_analysis.level }} RISK</div>
                      <div class="text-[11px] font-black uppercase tracking-[0.3em] mt-3 opacity-60">VERIFIED SIGNAL</div>
                    </div>

                    <div class="p-8 rounded-3xl bg-white/[0.03] border border-white/5 flex flex-col justify-between hover:border-primary/20 transition-all">
                      <div class="text-[10px] font-black text-gray-500 uppercase tracking-[0.4em] mb-4">Capital Density</div>
                      <div class="text-2xl font-black text-white uppercase tracking-tight">HIGH DEPTH</div>
                      <div class="text-[11px] text-gray-600 font-black uppercase tracking-[0.3em] mt-3 italic">LOW SLIPPAGE</div>
                    </div>
                  </div>

                  <!-- Insights & Concerns -->
                  <div class="grid md:grid-cols-2 gap-12 pt-8">
                    <div class="space-y-8 p-6 rounded-[32px] bg-rose-500/5 border border-rose-500/10">
                      <h4 class="text-[11px] font-black text-rose-500 uppercase tracking-[0.4em] flex items-center gap-4">
                        <span class="w-2 h-2 bg-rose-500 rounded-full shadow-[0_0_15px_rgba(244,63,94,0.6)]"></span> Critical Risk Factors
                      </h4>
                      <ul class="space-y-5">
                        <li v-for="(item, i) in parsedAnalysis.risk_analysis.concerns" :key="i" class="text-sm text-gray-300 font-bold flex gap-4 group/li leading-relaxed">
                          <span class="text-rose-500 group-hover:scale-150 transition-transform text-xl leading-none">•</span> 
                          <span class="group-hover:text-white transition-colors">{{ item }}</span>
                        </li>
                      </ul>
                    </div>
                    <div class="space-y-8 p-6 rounded-[32px] bg-cyan-400/5 border border-cyan-400/10">
                      <h4 class="text-[11px] font-black text-cyan-400 uppercase tracking-[0.4em] flex items-center gap-4">
                        <span class="w-2 h-2 bg-cyan-400 rounded-full shadow-[0_0_15px_rgba(34,211,238,0.6)]"></span> Strategic Intelligence
                      </h4>
                      <ul class="space-y-5">
                        <li v-for="(item, i) in parsedAnalysis.insights" :key="i" class="text-sm text-gray-300 font-bold flex gap-4 group/li leading-relaxed">
                          <span class="text-cyan-400 group-hover:scale-150 transition-transform text-xl leading-none">◈</span> 
                          <span class="group-hover:text-white transition-colors">{{ item }}</span>
                        </li>
                      </ul>
                    </div>
                  </div>
                </div>
                
                <!-- State UI -->
                <div v-else-if="aiError" class="bg-rose-500/10 border border-rose-500/20 p-10 rounded-[40px] text-rose-400 flex items-center gap-6 shadow-2xl relative z-10">
                   <span class="text-4xl">⚠️</span>
                   <div>
                     <div class="font-black uppercase tracking-[0.4em] text-sm mb-2">Protocol Error Detected</div>
                     <div class="text-xs font-bold opacity-70 leading-relaxed">{{ aiError }}</div>
                   </div>
                </div>
                
                <div v-else-if="analysis && !parsedAnalysis" class="prose prose-invert prose-lg max-w-none text-gray-300 font-medium py-6 px-4 relative z-10 leading-relaxed">
                   <div v-html="renderMarkdown(analysis)"></div>
                </div>

                <div v-if="!analysis && !isAnalyzing" class="flex flex-col items-center justify-center py-28 text-center relative z-10">
                   <div class="relative mb-12">
                     <div class="absolute inset-0 bg-primary/20 blur-[80px] rounded-full scale-150 animate-pulse"></div>
                     <div class="w-32 h-32 rounded-[40px] bg-gradient-to-br from-primary via-indigo-600 to-indigo-900 flex items-center justify-center text-6xl shadow-2xl relative z-10 border border-white/20 group-hover:rotate-6 transition-transform duration-700">🤖</div>
                   </div>
                   <h5 class="text-2xl font-black text-white mb-4 tracking-[0.3em] uppercase">Compute Alpha Insight</h5>
                   <p class="text-gray-500 text-xs mb-12 max-w-sm uppercase tracking-[0.4em] font-black leading-loose">
                     Initiate full-spectrum AI synthesis across social networks, exchange orderbooks, and on-chain protocol health.
                   </p>
                   <button 
                     @click="analyzeToken(displayToken)"
                     class="group relative px-20 py-8 rounded-[32px] bg-white text-black font-black text-sm uppercase tracking-[0.5em] transition-all hover:scale-105 active:scale-95 shadow-[0_25px_60px_rgba(255,255,255,0.2)] overflow-hidden"
                   >
                     <div class="absolute inset-0 bg-gradient-to-r from-primary to-indigo-600 opacity-0 group-hover:opacity-20 transition-opacity"></div>
                     <span class="relative z-10">Launch Protocol</span>
                   </button>
                </div>
             </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import RatingBadge from '../components/RatingBadge.vue'
import ScoreRadar from '../components/ScoreRadar.vue'
import MetricCard from '../components/MetricCard.vue'
import SentimentGauge from '../components/SentimentGauge.vue'
import { formatCurrency, formatPrice, formatNumber } from '../utils/formatters'
import { formatCurrencyCompact } from '../services/api'
import { useTokens } from '../composables/useTokens'
import { useAIAnalysis } from '../composables/useAIAnalysis'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

const route = useRoute()
const { fetchTokenDetail } = useTokens()
const { analysis, isAnalyzing, error: aiError, analyzeToken, clearAnalysis } = useAIAnalysis()

const detailData = ref(null)
const isLoading = ref(false)

const renderMarkdown = (text) => {
  if (!text) return ''
  const html = marked.parse(text)
  return DOMPurify.sanitize(html)
}

// Parse AI JSON response
const parsedAnalysis = computed(() => {
  if (!analysis.value) return null
  try {
    // If it's already an object (shouldn't happen with current useAIAnalysis but good for safety)
    if (typeof analysis.value === 'object') return analysis.value
    // Handle potential markdown code blocks
    let cleanJson = analysis.value
    if (cleanJson.includes('```json')) {
      cleanJson = cleanJson.split('```json')[1].split('```')[0].trim()
    } else if (cleanJson.includes('```')) {
      cleanJson = cleanJson.split('```')[1].trim()
    }
    return JSON.parse(cleanJson)
  } catch (e) {
    console.warn('Failed to parse AI JSON:', e)
    return null // Fallback to raw markdown display
  }
})

const riskClass = computed(() => {
  if (!parsedAnalysis.value) return ''
  const level = parsedAnalysis.value.risk_analysis.level.toLowerCase()
  if (level.includes('low')) return 'bg-green-500/5 border-green-500/20 text-green-400'
  if (level.includes('high')) return 'bg-rose-500/5 border-rose-500/20 text-rose-400'
  return 'bg-amber-500/5 border-amber-500/20 text-amber-400'
})

// Merge list data with detailed data
const displayToken = computed(() => {
  return { ...detailData.value }
})

onMounted(async () => {
  clearAnalysis() // Clear previous analysis
  const tokenId = route.params.id
  
  if (tokenId) {
    isLoading.value = true
    try {
      const data = await fetchTokenDetail(tokenId)
      if (data) {
        detailData.value = data
      }
    } finally {
      isLoading.value = false
    }
  }
})
</script>

<style scoped>
.glass-bento {
  background: rgba(255, 255, 255, 0.03);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.05);
  border-radius: 32px;
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.glass-bento:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(59, 130, 246, 0.2);
  transform: translateY(-4px);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.4);
}

.animate-shimmer {
  background-size: 200% 100%;
  animation: shimmer 3s linear infinite;
}

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

.shine {
  box-shadow: 0 0 10px currentColor;
}

.shadow-glow {
  box-shadow: 0 0 15px currentColor;
}

.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: rgba(59, 130, 246, 0.3);
}
</style>
