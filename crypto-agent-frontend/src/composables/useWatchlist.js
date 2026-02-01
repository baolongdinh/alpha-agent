import { ref, watch, onMounted } from 'vue'

const favorites = ref(new Set())

export function useWatchlist() {
  // Load from localStorage
  const loadFavorites = () => {
    try {
      const stored = localStorage.getItem('crypto_agent_favorites')
      if (stored) {
        favorites.value = new Set(JSON.parse(stored))
      }
    } catch (e) {
      console.error('Error loading favorites:', e)
    }
  }

  // Save to localStorage
  const saveFavorites = () => {
    try {
      localStorage.setItem('crypto_agent_favorites', JSON.stringify([...favorites.value]))
    } catch (e) {
      console.error('Error saving favorites:', e)
    }
  }

  const toggleFavorite = (symbol) => {
    if (favorites.value.has(symbol)) {
      favorites.value.delete(symbol)
    } else {
      favorites.value.add(symbol)
    }
    saveFavorites()
    // Trigger reactivity for Set
    favorites.value = new Set(favorites.value)
  }

  const isFavorite = (symbol) => {
    return favorites.value.has(symbol)
  }

  onMounted(() => {
    loadFavorites()
  })

  return {
    favorites,
    toggleFavorite,
    isFavorite
  }
}
