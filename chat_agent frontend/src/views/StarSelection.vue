<template>
  <div class="star-selection">
    <h1>选择您想聊天的明星</h1>
    <div class="stars-grid">
      <div
        v-for="star in stars"
        :key="star.id"
        class="star-card"
        @click="selectStar(star.id)"
      >
        <div class="star-avatar">
          {{ star?.name ? star.name.charAt(0) : '?' }}
        </div>
        <h3>{{ star?.name || '未知明星' }}</h3>
        <p>{{ star?.introduction || '暂无介绍' }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()
const stars = ref([])

onMounted(async () => {
  try {
    // 实际调用后端API获取明星列表
    const response = await axios.get('http://localhost:8000/api/v1/stars')
    stars.value = response.data.data || []
  } catch (error) {
    console.error('获取明星列表失败:', error)
    // 错误时使用空数组避免UI问题
    stars.value = []
  }
})

const selectStar = (starId) => {
  router.push({ name: 'ChatPage', params: { starId } })
}
</script>

<style scoped>
.star-selection {
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem;
}

.star-selection h1 {
  text-align: center;
  margin-bottom: 3rem;
  color: #2c3e50;
  font-size: 2.5rem;
  font-weight: 700;
  letter-spacing: -0.5px;
}

.stars-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  grid-auto-flow: row;
  gap: 2rem;
  justify-items: center;
  width: 100%;
  
  /* 响应式设计 */
  @media (max-width: 1200px) {
    grid-template-columns: repeat(3, 1fr);
  }
  @media (max-width: 900px) {
    grid-template-columns: repeat(2, 1fr);
  }
  @media (max-width: 600px) {
    grid-template-columns: 1fr;
  }
}

.star-card {
  background: #fff;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  cursor: pointer;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  text-align: center;
  width: 100%;
  max-width: 300px;
}

.star-card:hover {
  transform: translateY(-8px) scale(1.02);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.15);
}

.star-avatar {
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
  margin: 0 auto 1rem;
}

.star-card h3 {
  margin-bottom: 0.5rem;
  color: #333;
}

.star-card p {
  color: #666;
  font-size: 0.9rem;
}
</style>