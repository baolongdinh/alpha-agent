<template>
  <div class="min-h-screen bg-[#050510]">
    <!-- Full Page Container -->
    <div class="p-0 overflow-hidden bg-[#0a0b0f] w-full mx-auto flex flex-col min-h-screen relative text-gray-300">
      
      <!-- HERO HEADER SECTION -->
      <div class="relative p-6 pt-24 pb-12 overflow-hidden">
        <!-- Dynamic Hero Background with sophisticated mesh -->
        <div class="absolute inset-0 bg-[#0a0b0f] z-0"></div>
        <div class="absolute inset-0 bg-[radial-gradient(circle_at_20%_30%,rgba(59,130,246,0.15)_0%,transparent_50%)] z-0"></div>
        <div class="absolute inset-0 bg-[radial-gradient(circle_at_80%_70%,rgba(99,102,241,0.1)_0%,transparent_50%)] z-0"></div>
        <div class="absolute top-0 right-0 w-[800px] h-[800px] bg-primary/20 blur-[150px] rounded-full translate-x-1/3 -translate-y-1/3 z-0 animate-slow-pulse"></div>
        
        <!-- Redesigned Back Button: Elegant Translucent Pill -->
        <button 
          @click="$router.push('/')" 
          class="absolute top-8 left-8 flex items-center gap-3 px-5 py-2.5 rounded-full bg-white/[0.03] hover:bg-white/[0.08] border border-white/10 backdrop-blur-xl transition-all duration-300 group z-50 hover:shadow-[0_0_20px_rgba(255,255,255,0.05)]"
        >
          <div class="w-7 h-7 rounded-full bg-white/5 flex items-center justify-center group-hover:bg-primary/20 transition-colors">
            <svg class="w-4 h-4 text-gray-400 group-hover:text-primary transition-colors duration-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M15 19l-7-7 7-7"/>
            </svg>
          </div>
          <span class="text-[10px] font-black uppercase tracking-[0.2em] text-gray-400 group-hover:text-white transition-colors">Return Pulse</span>
        </button>

        <div v-if="displayToken && displayToken.symbol" class="relative z-10 max-w-[1400px] mx-auto">
          <div class="flex flex-col lg:flex-row lg:items-end justify-between gap-10">
            <!-- Left: Identity with enhanced depth -->
            <div class="flex items-center gap-8">
              <div class="relative group">
                <div class="absolute -inset-6 bg-primary/25 blur-3xl rounded-full opacity-40 group-hover:opacity-100 transition-opacity duration-1000"></div>
                <div class="w-24 h-24 rounded-[32px] bg-gradient-to-br from-white/10 to-transparent flex items-center justify-center border border-white/20 shadow-2xl relative z-10 overflow-hidden backdrop-blur-3xl group-hover:scale-105 transition-transform duration-500">
                  <img v-if="displayToken.image" :src="displayToken.image" :alt="displayToken.name" class="w-full h-full object-cover p-3 group-hover:p-1 transition-all duration-700" />
                  <span v-else class="text-4xl font-black text-primary">{{ displayToken.symbol?.[0] || '?' }}</span>
                </div>
              </div>
              
              <div class="space-y-3">
                <div class="flex items-center gap-4 flex-wrap">
                  <h2 class="text-5xl md:text-6xl font-black text-white tracking-tighter leading-none">{{ displayToken.name }}</h2>
                  <div class="flex items-center gap-2">
                    <span class="px-3 py-1 bg-white/5 border border-white/10 rounded-lg text-gray-400 text-[10px] font-mono font-black tracking-widest uppercase">{{ displayToken.symbol }}</span>
                    <div class="px-3 py-1 bg-primary/10 text-primary text-[10px] font-black rounded-lg border border-primary/20 flex items-center gap-2 shadow-inner">
                      <span class="w-1.5 h-1.5 bg-primary rounded-full animate-ping"></span>
                      RANK #{{ displayToken.rank }}
                    </div>
                  </div>
                </div>
                <div class="flex items-center gap-2">
                  <div v-if="displayToken.category" class="px-3 py-1 bg-white/5 text-gray-500 text-[9px] font-black rounded-md border border-white/5 uppercase tracking-widest hover:text-white transition-colors cursor-default">
                    {{ displayToken.category }}
                  </div>
                  <div class="px-3 py-1 bg-white/5 text-gray-500 text-[9px] font-black rounded-md border border-white/5 uppercase tracking-widest">
                    Live Analysis
                  </div>
                </div>
              </div>
            </div>

            <!-- Right: Premium Real-time Price -->
            <div class="flex flex-col items-end gap-2 text-right">
              <div class="text-[10px] font-black text-primary uppercase tracking-[0.5em] mb-1 opacity-70">Market Valuation</div>
              <div class="text-6xl font-black text-white tracking-tighter mb-2 tabular-nums drop-shadow-[0_0_30px_rgba(255,255,255,0.1)]">
                {{ formatCurrency(displayToken.price) }}
              </div>
              
              <div class="flex items-center gap-4 justify-end">
                <div 
                  class="flex items-center gap-2 px-4 py-2 rounded-2xl font-black text-xs border backdrop-blur-md shadow-xl transition-transform hover:scale-110"
                  :class="displayToken.change_24h >= 0 ? 'bg-green-500/10 border-green-500/20 text-green-400' : 'bg-red-500/10 border-red-500/20 text-red-500'"
                >
                  <span class="text-xl leading-none">{{ displayToken.change_24h >= 0 ? '↑' : '↓' }}</span>
                  <span class="font-mono text-sm">{{ Math.abs(displayToken.change_24h).toFixed(2) }}%</span>
                  <span class="text-[9px] opacity-50 uppercase font-black">24H</span>
                </div>
                
                <div 
                  class="flex items-center gap-2 px-4 py-2 rounded-2xl font-black text-xs border backdrop-blur-md shadow-xl transition-transform hover:scale-110"
                  :class="displayToken.change_7d >= 0 ? 'bg-green-500/10 border-green-500/10 text-green-400/70' : 'bg-red-500/10 border-red-500/10 text-red-500/70'"
                >
                   <span class="font-mono text-sm">{{ Math.abs(displayToken.change_7d).toFixed(2) }}%</span>
                   <span class="text-[9px] opacity-50 uppercase font-black">7D</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Refined Quick Stats Belt -->
          <div class="grid grid-cols-2 md:grid-cols-4 gap-8 mt-12 pt-8 border-t border-white/10 relative z-10">
            <div class="group py-2">
              <div class="text-[10px] text-gray-500 font-black uppercase tracking-[0.3em] mb-3 group-hover:text-primary transition-colors flex items-center gap-2">
                <span class="w-1.5 h-1.5 bg-primary/40 rounded-full group-hover:bg-primary"></span>
                Circulating
              </div>
              <div class="text-xl font-mono text-white flex items-end gap-2 group-hover:scale-105 transition-transform origin-left">
                {{ formatNumber(displayToken.circulating_supply) }} 
                <span class="text-[10px] text-gray-600 mb-1 font-sans font-black uppercase tracking-widest">{{ displayToken.symbol }}</span>
              </div>
              <div class="w-full h-1 bg-white/5 rounded-full mt-3 overflow-hidden border border-white/5">
                <div class="h-full bg-gradient-to-r from-primary to-blue-400 shadow-[0_0_10px_rgba(59,130,246,0.5)]" :style="{ width: displayToken.max_supply > 0 ? (displayToken.circulating_supply / displayToken.max_supply * 100) + '%' : '100%' }"></div>
              </div>
            </div>
            <div class="group py-2">
              <div class="text-[10px] text-gray-500 font-black uppercase tracking-[0.3em] mb-3 group-hover:text-amber-500 transition-colors flex items-center gap-2">
                <span class="w-1.5 h-1.5 bg-amber-500/40 rounded-full group-hover:bg-amber-500"></span>
                Max Supply
              </div>
              <div class="text-xl font-mono text-white group-hover:scale-105 transition-transform origin-left">{{ displayToken.max_supply > 0 ? formatNumber(displayToken.max_supply) : '∞ UNLIMITED' }}</div>
            </div>
            <div class="group py-2">
              <div class="text-[10px] text-gray-500 font-black uppercase tracking-[0.3em] mb-3 group-hover:text-cyan-500 transition-colors flex items-center gap-2">
                <span class="w-1.5 h-1.5 bg-cyan-500/40 rounded-full group-hover:bg-cyan-500"></span>
                Market Share
              </div>
              <div class="text-xl font-mono text-cyan-400 font-black group-hover:scale-110 transition-transform origin-left">{{ (displayToken.MarketCapDom || displayToken.market_cap_dominance || 0).toFixed(3) }}%</div>
            </div>
            <div class="group py-2">
              <div class="text-[10px] text-gray-500 font-black uppercase tracking-[0.3em] mb-3 group-hover:text-green-500 transition-colors flex items-center gap-2">
                <span class="w-1.5 h-1.5 bg-green-500/40 rounded-full group-hover:bg-green-500"></span>
                Valuation (FDV)
              </div>
              <div class="text-xl font-mono text-white tracking-tight group-hover:scale-105 transition-transform origin-left">{{ formatCurrencyCompact(displayToken.fdv) }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- MAIN APP BODY -->
      <div class="flex-1 overflow-y-auto custom-scrollbar bg-[#050608] p-6 rounded-t-[40px] border-t border-white/10 relative z-10 -mt-6 shadow-[0_-50px_100px_rgba(0,0,0,0.5)]">
        <div class="grid lg:grid-cols-12 gap-6">
          
          <!-- LEFT COLUMN: Advanced Scoring & Metrics (4 cols) -->
          <div class="lg:col-span-4 space-y-6">
            <!-- Global Alpha Score Card -->
            <div class="glass-bento p-6 flex flex-col items-center group relative overflow-hidden">
              <div class="absolute -top-10 -left-10 w-40 h-40 bg-primary/10 blur-[60px] rounded-full group-hover:bg-primary/20 transition-all duration-700"></div>
              <div class="text-[9px] font-black text-gray-500 uppercase tracking-[0.4em] mb-4 w-full text-center">Core Trust Rating</div>
              
              <div class="relative mb-4">
                <div class="absolute inset-0 bg-primary/20 blur-3xl rounded-full scale-150 opacity-40"></div>
                <RatingBadge 
                  :grade="displayToken.score_breakdown?.grade || 'D'" 
                  class="text-4xl px-8 py-4 rounded-[24px] shadow-[0_20px_60px_rgba(0,0,0,0.5)] border-2 border-white/10 relative z-10"
                />
              </div>
              
              <div class="text-center group-hover:scale-105 transition-transform">
                <div class="text-5xl font-black text-white tracking-tighter leading-none">{{ displayToken.trust_score?.toFixed(0) || '—' }}</div>
                <div class="text-[10px] text-primary font-black uppercase tracking-[0.3em] mt-2">Alpha Index</div>
              </div>
            </div>

            <!-- Sentiment & Vibe -->
            <div class="glass-bento p-6 text-center text-left">
              <div class="text-[9px] font-black text-gray-500 uppercase tracking-[0.4em] mb-6 text-left">Community Vibe</div>
              <div class="flex flex-col items-center justify-center py-2">
                <SentimentGauge :score="displayToken.sentiment !== undefined ? displayToken.sentiment : ((displayToken.trust_score || 50) / 100)" class="scale-100" />
                <div class="mt-6 text-[10px] font-black text-gray-400 bg-white/5 px-4 py-2 rounded-xl border border-white/5 flex items-center gap-2">
                  <span class="w-1.5 h-1.5 rounded-full" :class="(displayToken.sentiment !== undefined ? displayToken.sentiment : ((displayToken.trust_score || 50) / 100)) >= 0.5 ? 'bg-green-400 shine' : 'bg-rose-400 shine'"></span>
                  SENTIMENT SHIFT: <span :class="(displayToken.sentiment !== undefined ? displayToken.sentiment : ((displayToken.trust_score || 50) / 100)) >= 0.5 ? 'text-green-400' : 'text-rose-400'">
                    {{ ((displayToken.sentiment !== undefined ? displayToken.sentiment : ((displayToken.trust_score || 50) / 100)) * 100).toFixed(1) }}% MATCH
                  </span>
                </div>
              </div>
            </div>

            <!-- Radar Matrix (Scaled Down) -->
            <div class="glass-bento p-6 flex flex-col items-center justify-center min-h-[350px]">
              <div class="text-[9px] font-black text-gray-500 uppercase tracking-[0.4em] mb-6 w-full text-left">Asset Signature Matrix</div>
              <div class="w-full h-[280px] relative">
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
                   label="Monthly Performance" 
                   :value="displayToken.change_30d ? (displayToken.change_30d >= 0 ? '+' : '') + displayToken.change_30d.toFixed(1) + '%' : '—'" 
                   :change="displayToken.change_30d"
                   :isPositive="displayToken.change_30d >= 0"
                   icon="�"
                   subLabel="30-Day Trajectory"
                   :progress="Math.min(Math.abs(displayToken.change_30d || 0), 100)"
                   :progressColor="displayToken.change_30d >= 0 ? 'green' : 'rose'"
                />
              </div>
            </div>

            <!-- Section: Advanced Market Dynamics -->
            <div class="space-y-4">
              <h4 class="text-[10px] font-black text-amber-500/80 uppercase tracking-[0.4em] px-4">Market Dynamics</h4>
              <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                <!-- 1. SENTIMENT SHIFT (Replaces Holders) -->
                <MetricCard 
                   label="Sentiment Shift" 
                   :value="((displayToken.trust_score || 50) > 60) ? 'BULLISH' : ((displayToken.trust_score || 50) < 40 ? 'BEARISH' : 'NEUTRAL')" 
                   icon="🧠"
                   :subLabel="`Confidence Score: ${displayToken.trust_score?.toFixed(0) || 50}/100`"
                   class="hover:border-amber-500/30 transition-all duration-500"
                   :progress="displayToken.trust_score || 50"
                   progressColor="amber"
                />
                
                <!-- 2. LIQUIDITY TURNOVER (Replaces Liquidity) -->
                <MetricCard 
                   label="Liquidity Turnover" 
                   :value="displayToken.market_cap > 0 ? ((displayToken.volume_24h / displayToken.market_cap) * 100).toFixed(2) + '%' : '—'" 
                   icon="🌊"
                   :subLabel="((displayToken.volume_24h / displayToken.market_cap) > 0.1) ? 'High Velocity' : 'Standard Velocity'"
                   :progress="displayToken.market_cap > 0 ? Math.min((displayToken.volume_24h / displayToken.market_cap) * 500, 100) : 0"
                   progressColor="cyan"
                />

                <!-- 3. VOLATILITY INDEX (Replaces Revenue) -->
                <MetricCard 
                   label="Volatility Index" 
                   :value="Math.abs(displayToken.change_30d || displayToken.change_7d || 0).toFixed(1) + '%'" 
                   icon="📊"
                   subLabel="30-Day Deviation"
                   :progress="Math.min(Math.abs(displayToken.change_30d || displayToken.change_7d || 0), 100)"
                   progressColor="emerald"
                />
              </div>
            </div>

            <!-- AI Alpha Agent Insight -->
            <div class="relative group mt-2">
              <div class="absolute -inset-1 bg-gradient-to-r from-primary/30 via-indigo-500/20 to-primary/30 rounded-[30px] blur-xl opacity-40 group-hover:opacity-100 transition duration-1000"></div>
              
              <div class="glass-bento p-8 relative overflow-hidden bg-black/80 backdrop-blur-3xl border border-white/10 rounded-[30px] shadow-2xl">
                <!-- Decorative Elements -->
                <div class="absolute top-0 right-0 w-full h-full bg-[radial-gradient(circle_at_100%_0%,rgba(59,130,246,0.1)_0%,transparent_50%)]"></div>
                <div class="absolute -bottom-24 -left-24 w-96 h-96 bg-primary/5 blur-[120px] rounded-full"></div>
                
                <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 mb-8 relative z-10">
                   <div class="flex items-center gap-4">
                     <div class="w-12 h-12 rounded-[18px] bg-primary/10 flex items-center justify-center text-primary text-2xl shadow-lg border border-primary/20 group-hover:scale-110 transition-transform duration-500">✨</div>
                     <div>
                       <h3 class="text-xl font-black text-white tracking-widest uppercase">Agent Alpha Report</h3>
                       <div class="text-[9px] text-primary font-black uppercase tracking-[0.4em] flex items-center gap-1.5 mt-1.5">
                          <span class="w-1.5 h-1.5 bg-primary rounded-full animate-ping"></span>
                          Full Spectrum Synthesis Active
                       </div>
                     </div>
                   </div>
                   <div v-if="isAnalyzing" class="flex items-center gap-3 px-4 py-2 bg-primary/10 border border-primary/20 rounded-xl backdrop-blur-md">
                     <div class="flex gap-1.5">
                        <span class="w-1.5 h-3 bg-primary rounded-full animate-bounce [animation-delay:-0.3s]"></span>
                        <span class="w-1.5 h-3 bg-primary rounded-full animate-bounce [animation-delay:-0.15s]"></span>
                        <span class="w-1.5 h-3 bg-primary rounded-full animate-bounce"></span>
                     </div>
                     <span class="text-[10px] font-black text-primary uppercase tracking-[0.2em] animate-pulse">Processing Market Signals...</span>
                   </div>
                </div>
                
                <!-- Structured JSON AI UI -->
                <div v-if="parsedAnalysis" class="space-y-8 relative z-10">
                  <!-- Executive Summary -->
                  <div class="relative">
                    <div class="absolute -left-6 top-0 bottom-0 w-1 bg-gradient-to-b from-primary to-transparent rounded-full opacity-60"></div>
                    <p class="text-xl md:text-2xl font-black text-white leading-tight tracking-tight pl-6 pr-4 italic">
                      "{{ parsedAnalysis.summary }}"
                    </p>
                  </div>

                  <div class="grid md:grid-cols-2 gap-4">
                    <!-- Potential Score UI -->
                    <div class="p-6 rounded-2xl bg-white/[0.03] border border-white/5 space-y-4 group/potential hover:bg-white/[0.05] transition-all">
                      <div class="flex justify-between items-end">
                         <span class="text-[9px] font-black text-gray-500 uppercase tracking-[0.3em]">Future Growth Potential</span>
                         <span class="text-4xl font-black text-primary tabular-nums drop-shadow-[0_0_15px_rgba(59,130,246,0.5)]">
                           {{ parsedAnalysis.growth_potential.score }}<span class="text-base opacity-50 ml-1 font-bold">%</span>
                         </span>
                      </div>
                      <div class="h-1.5 w-full bg-black/60 rounded-full overflow-hidden border border-white/10">
                        <div class="h-full rounded-full bg-gradient-to-r from-primary via-indigo-400 to-primary bg-[length:200%_100%] animate-shimmer transition-all duration-1500" :style="{ width: `${parsedAnalysis.growth_potential.score}%` }"></div>
                      </div>
                      <p class="text-[11px] text-gray-400 leading-relaxed font-medium italic mt-2">{{ parsedAnalysis.growth_potential.reason }}</p>
                    </div>

                    <!-- Recommendation Card -->
                    <div class="p-6 rounded-2xl bg-gradient-to-br from-primary/20 via-primary/5 to-transparent border border-primary/20 shadow-xl transition-all group/recommend hover:border-primary/40">
                      <div class="text-[9px] font-black text-primary uppercase tracking-[0.3em] mb-4">Strategic Recommendation</div>
                      <div class="flex items-center justify-between mb-4">
                         <span class="text-3xl font-black text-white tracking-widest uppercase">{{ parsedAnalysis.recommendation.action }}</span>
                         <div class="text-[9px] font-black text-primary bg-primary/20 px-3 py-1.5 rounded-lg border border-primary/30 tracking-[0.1em] font-mono shadow-lg">ZONE: {{ parsedAnalysis.recommendation.entry_zone }}</div>
                      </div>
                      <div class="pt-4 border-t border-white/10 text-[9px] text-gray-500 font-bold uppercase tracking-[0.2em]">Target Exit: <span class="text-white ml-2 text-xs font-black">{{ parsedAnalysis.recommendation.target }}</span></div>
                    </div>
                  </div>

                  <!-- Execution Plan (New) -->
                  <div v-if="parsedAnalysis.trading_plan" class="glass-bento p-6 relative overflow-hidden group/plan">
                    <div class="absolute top-0 right-0 w-64 h-64 bg-green-500/5 blur-[80px] rounded-full group-hover/plan:bg-green-500/10 transition-all duration-700"></div>
                    <h4 class="text-[9px] font-black text-green-400 uppercase tracking-[0.3em] mb-4 relative z-10 flex items-center gap-2">
                      <span class="w-1.5 h-1.5 rounded-full bg-green-400 shadow-glow"></span> Strategic Execution Plan
                    </h4>
                    
                    <div class="grid md:grid-cols-4 gap-4 relative z-10 mb-6">
                        <!-- Entry -->
                        <div class="col-span-2 p-4 rounded-xl bg-white/5 border border-white/5 hover:border-green-500/30 transition-colors">
                            <div class="text-[9px] text-gray-500 font-bold uppercase tracking-widest mb-2 italic">Accumulation Strategy</div>
                            <div class="text-[11px] font-bold text-gray-200 leading-relaxed">{{ parsedAnalysis.trading_plan.buy_strategy }}</div>
                        </div>

                        <!-- Stop Loss -->
                        <div class="p-4 rounded-xl bg-white/5 border border-white/5 hover:border-rose-500/30 transition-colors flex flex-col justify-center">
                             <div class="text-[9px] text-gray-500 font-bold uppercase tracking-widest mb-1 italic">Invalidation Level</div>
                             <div class="text-base font-black text-rose-400">{{ parsedAnalysis.trading_plan.stop_loss }}</div>
                        </div>

                        <!-- Timeframe -->
                        <div class="p-4 rounded-xl bg-white/5 border border-white/5 flex flex-col justify-center">
                             <div class="text-[9px] text-gray-500 font-bold uppercase tracking-widest mb-1 italic">Horizon</div>
                             <div class="text-base font-black text-white uppercase tracking-tight">{{ parsedAnalysis.trading_plan.time_horizon }}</div>
                        </div>
                    </div>

                    <!-- Take Profit Targets -->
                    <div class="relative z-10">
                        <div class="text-[9px] text-gray-500 font-bold uppercase tracking-widest mb-3 italic">Profit Realization Targets</div>
                        <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
                            <div v-for="(target, idx) in parsedAnalysis.trading_plan.sell_targets" :key="idx" 
                                 class="px-4 py-2 rounded-xl bg-gradient-to-r from-green-500/10 to-transparent border border-green-500/20 flex items-center gap-3 hover:bg-green-500/20 transition-all">
                                <span class="w-5 h-5 rounded-full bg-green-500/20 flex items-center justify-center text-green-400 font-black text-[9px] border border-green-500/30">{{ idx + 1 }}</span>
                                <span class="font-bold text-green-100 text-xs">{{ target }}</span>
                            </div>
                        </div>
                    </div>
                  </div>

                  <!-- SECTION: Fundamental Intelligence -->
                  <div v-if="parsedAnalysis.fundamental_analysis" class="mb-2 animate-slide-up-delay-2">
                    <div class="glass-bento p-6 relative overflow-hidden">
                       <h3 class="text-[9px] font-black uppercase tracking-[0.3em] text-purple-400 mb-6 flex items-center gap-2">
                          <span class="w-1.5 h-1.5 bg-purple-400 rounded-full animate-pulse shadow-glow"></span>
                          Fundamental Intelligence
                       </h3>

                       <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                          <!-- Sector -->
                          <div class="bg-white/5 rounded-xl p-4 border border-white/5 hover:border-purple-500/30 transition-colors">
                             <div class="text-[9px] uppercase tracking-[0.2em] text-gray-500 mb-2 font-bold italic">Core Sector</div>
                             <div class="text-gray-200 font-bold text-xs leading-relaxed">
                                {{ parsedAnalysis.fundamental_analysis.sector }}
                             </div>
                          </div>

                          <!-- Tokenomics -->
                          <div class="bg-white/5 rounded-xl p-4 border border-white/5 hover:border-purple-500/30 transition-colors">
                             <div class="text-[9px] uppercase tracking-[0.2em] text-gray-500 mb-2 font-bold italic">Tokenomics Structure</div>
                             <div class="text-gray-200 font-bold text-xs leading-relaxed">
                                {{ parsedAnalysis.fundamental_analysis.tokenomics }}
                             </div>
                          </div>

                          <!-- Moat -->
                          <div class="bg-white/5 rounded-xl p-4 border border-white/5 hover:border-purple-500/30 transition-colors">
                             <div class="text-[9px] uppercase tracking-[0.2em] text-gray-500 mb-2 font-bold italic">Economic Moat</div>
                             <div class="text-gray-200 font-bold text-xs leading-relaxed">
                                {{ parsedAnalysis.fundamental_analysis.economic_moat }}
                             </div>
                          </div>
                       </div>
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
                     <span class="relative z-10">Analyze With AI</span>
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
@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

@keyframes slow-pulse {
  0%, 100% { opacity: 0.4; transform: scale(1) translate(33%, -33%); }
  50% { opacity: 0.6; transform: scale(1.1) translate(30%, -30%); }
}

.animate-slow-pulse {
  animation: slow-pulse 10s ease-in-out infinite;
}

.glass-bento {
  background: rgba(255, 255, 255, 0.02);
  backdrop-filter: blur(24px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.04);
  border-radius: 40px;
  transition: all 0.6s cubic-bezier(0.34, 1.56, 0.64, 1);
  box-shadow: 
    0 10px 30px -10px rgba(0, 0, 0, 0.5),
    inset 0 1px 1px rgba(255, 255, 255, 0.05);
}

.glass-bento:hover {
  background: rgba(255, 255, 255, 0.04);
  border-color: rgba(59, 130, 246, 0.3);
  transform: translateY(-8px) scale(1.01);
  box-shadow: 
    0 30px 60px -20px rgba(0, 0, 0, 0.7),
    0 0 20px rgba(59, 130, 246, 0.1);
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
