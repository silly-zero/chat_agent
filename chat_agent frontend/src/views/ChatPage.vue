<template>
  <div class="chat-page">
    <!-- 侧边栏 -->
    <aside :class="sidebarClass">
      <div class="sidebar-header">
        <button @click="goBack" class="back-button-mobile">←</button>
        <h2>对话</h2>
      </div>
      
      <div class="star-info">
        <div class="star-avatar-large">
          {{ currentStar?.name?.charAt(0) || '?' }}
        </div>
        <div class="star-details">
          <h3>{{ currentStar?.name || '加载中...' }}</h3>
          <p>{{ currentStar?.introduction || '' }}</p>
        </div>
      </div>
      
      <div class="conversation-memories">
        <h3>对话记录</h3>
        <div class="memory-list">
          <div
            v-for="memory in conversationMemories"
            :key="memory.id"
            class="memory-item"
            @click="loadMemory(memory)"
          >
            <div class="memory-preview">{{ memory.preview }}</div>
            <div class="memory-time">{{ memory.timestamp }}</div>
          </div>
        </div>
      </div>
    </aside>
    
    <!-- 主聊天区域 -->
    <main class="chat-main">
      <header class="chat-header">
        <button @click="goBack" class="back-button">←</button>
        <div class="header-info">
          <h2>{{ currentStar?.name || '加载中...' }}</h2>
        </div>
        <div class="header-actions">
          <button class="new-chat-button" @click="startNewChat">新对话</button>
        </div>
      </header>
      
      <div class="chat-history-container">
        <div class="chat-history" ref="chatHistory">
          <!-- 欢迎消息 -->
          <div v-if="messages.length === 0" class="welcome-message">
            <p>开始与 {{ currentStar?.name || '明星' }} 聊天吧！</p>
          </div>
          
          <!-- 消息列表 -->
          <div v-for="message in messages" :key="message.id" :class="['message', message.role]">
            <div v-if="message.role === 'bot'" class="message-avatar">
              {{ currentStar?.name?.charAt(0) || '?' }}
            </div>
            <div v-if="message.role === 'user'" class="message-avatar user-avatar">
              <span>我</span>
            </div>
            <div class="message-wrapper">
              <div class="message-content">
                {{ message.content }}
              </div>
              <div class="message-time">{{ message.timestamp }}</div>
            </div>
          </div>
          
          <!-- 加载状态 -->
          <div v-if="isLoading" class="message bot loading">
            <div class="message-avatar">
              {{ currentStar?.name?.charAt(0) || '?' }}
            </div>
            <div class="message-content">
              <span class="loading-dot"></span>
              <span class="loading-dot"></span>
              <span class="loading-dot"></span>
            </div>
          </div>
        </div>
      </div>
      
      <!-- 输入区域 -->
      <div class="chat-input-container">
        <div class="input-wrapper">
          <textarea
            v-model="inputMessage"
            @keydown.enter.ctrl="sendMessage"
            @keydown.enter.meta="sendMessage"
            placeholder="发消息或输入"
            class="chat-input"
            :disabled="isLoading"
            rows="3"
          ></textarea>
          <div class="input-right-actions">
            <button @click="sendMessage" class="send-button-icon" :disabled="isLoading || !inputMessage.trim()">
              <span class="icon">↑</span>
            </button>
          </div>
        </div>
      </div>
    </main>
    
    <!-- 移动端侧边栏切换按钮 -->
    <button class="sidebar-toggle" @click="toggleSidebar">
      ☰
    </button>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

const route = useRoute()
const router = useRouter()
const chatHistory = ref(null)

const starId = ref(route.params.starId)
const currentStar = ref(null)
const messages = ref([])
const inputMessage = ref('')
const isLoading = ref(false)
const conversationMemories = ref([])
const sidebarVisible = ref(false)
const currentChatId = ref(null)

onMounted(() => {
  loadStarInfo()
  loadConversationMemories()
})

const loadStarInfo = async () => {
  try {
    // 实际调用后端API获取明星信息
    const response = await axios.get(`http://localhost:8000/api/v1/stars/${starId.value}`)
    currentStar.value = response.data.data
  } catch (error) {
    console.error('获取明星信息失败:', error)
  }
}

const loadConversationMemories = async () => {
  try {
    // 实际调用后端API获取对话记忆
    const response = await axios.get(`http://localhost:8000/api/v1/chats?star_id=${starId.value}`)
    conversationMemories.value = (response.data.data || []).map(chat => ({
      id: chat.id,
      preview: chat.last_message || '无消息',
      timestamp: new Date(chat.last_active).toLocaleString('zh-CN'),
      messages: [] // 消息将在加载时获取
    }))
  } catch (error) {
    console.error('获取对话记忆失败:', error)
    // 错误时使用空数组避免UI问题
    conversationMemories.value = []
  }
}


  
  const sendMessage = async () => {
    if (!inputMessage.value.trim() || isLoading.value) return
    
    const currentTime = new Date();
    const userMessage = {
      id: Date.now(),
      role: 'user',
      content: inputMessage.value.trim(),
      timestamp: currentTime.toLocaleString('zh-CN'),
      createdAt: currentTime.getTime() // 添加时间戳用于排序
    }
    
    messages.value.push(userMessage);
    
    // 按照时间戳排序以确保顺序正确
    messages.value.sort((a, b) => (a.createdAt || 0) - (b.createdAt || 0));
    
    // 移除临时的createdAt字段
    messages.value = messages.value.map(({ createdAt, ...rest }) => rest);
    inputMessage.value = ''
    scrollToBottom()
    
    isLoading.value = true
    
    try {
      // 先创建或获取聊天会话
      let chatId = currentChatId.value;
      
      if (!chatId) {
        // 确保starId是数字类型
        const numericStarId = Number(starId.value);
        
        // 使用fetch API发送请求
        const response = await fetch('http://localhost:8000/api/v1/chats', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ star_id: numericStarId })
        });
        
        // 获取响应文本
        const responseText = await response.text();
        
        if (!response.ok) {
          throw new Error(`创建聊天会话失败: ${response.status} ${responseText}`);
        }
        
        // 解析响应数据
        const data = JSON.parse(responseText);
        
        // 从响应中获取chatId - 兼容多种响应格式
        if (data && data.code === 200 && data.data && data.data.id) {
          chatId = data.data.id;
        } else if (data && data.id) {
          // 兼容直接返回id的格式
          chatId = data.id;
        } else {
          throw new Error(`响应格式异常，未找到chatId: ${JSON.stringify(data)}`);
        }
        
        currentChatId.value = chatId;
      }
      
      // 检查chatId是否有效
      if (!chatId || (typeof chatId !== 'number' && typeof chatId !== 'string')) {
        throw new Error('无效的聊天ID: ' + chatId);
      }
      
      // 发送消息
      const messageRequestData = {
        chat_id: chatId,
        content: userMessage.content,
        model: 'doubao-1.5-pro-32k-250115' // 使用豆包模型
      };
      
      const response = await axios.post('http://localhost:8000/api/v1/chats/messages', messageRequestData);
      
      // 处理后端返回的响应格式
      const responseData = response.data;
      
      // 从响应数据中正确提取消息内容
      let messageContent = '消息内容为空';
      let messageId = Date.now() + 1;
      
      if (responseData && responseData.code === 200 && responseData.data && responseData.data.content) {
        // 正常情况：从 responseData.data.content 获取内容
        messageContent = responseData.data.content;
        messageId = responseData.data.id || messageId;
      } else if (responseData && responseData.content) {
        // 兼容情况：直接从 responseData.content 获取
        messageContent = responseData.content;
        messageId = responseData.id || messageId;
      } else {
        throw new Error('无法提取消息内容');
      }
      
      const currentTime = new Date();
      const botMessage = {
        id: messageId,
        role: 'bot',
        content: messageContent,
        timestamp: currentTime.toLocaleString('zh-CN'),
        createdAt: currentTime.getTime() // 添加时间戳用于排序
      }
      
      messages.value.push(botMessage);
      
      // 按照时间戳排序以确保顺序正确
      messages.value.sort((a, b) => (a.createdAt || 0) - (b.createdAt || 0));
      
      // 移除临时的createdAt字段
      messages.value = messages.value.map(({ createdAt, ...rest }) => rest);
      isLoading.value = false
      scrollToBottom()
      
      // 保存对话到记忆
      saveConversationToMemory()
    } catch (error) {
      isLoading.value = false
      
      // 简化错误处理
      let errorDetail = '';
      if (error.response) {
        errorDetail = `服务器错误: ${error.response.status}`;
      } else if (error.request) {
        errorDetail = '服务器无响应，请检查网络连接';
      } else {
        errorDetail = error.message || '未知错误';
      }
      
      // 显示错误提示
      const currentTime = new Date();
      const errorMessage = {
        id: Date.now() + 1,
        role: 'bot',
        content: `抱歉，消息发送失败，请稍后重试。\n错误详情: ${errorDetail}`,
        timestamp: currentTime.toLocaleString('zh-CN'),
        createdAt: currentTime.getTime() // 添加时间戳用于排序
      }
      messages.value.push(errorMessage);
      
      // 按照时间戳排序以确保顺序正确
      messages.value.sort((a, b) => (a.createdAt || 0) - (b.createdAt || 0));
      
      // 移除临时的createdAt字段
      messages.value = messages.value.map(({ createdAt, ...rest }) => rest);
      scrollToBottom()
    }
}

const saveConversationToMemory = async () => {
  // 保存最近的对话到记忆
  if (messages.value.length > 0 && currentChatId.value) {
    try {
      // 更新聊天会话信息
      await axios.put(`http://localhost:8000/api/v1/chats/${currentChatId.value}`, {
        title: messages.value[0]?.content?.substring(0, 30) + '...' || '新对话'
      })
      
      // 更新本地对话记忆列表
      await loadConversationMemories()
    } catch (error) {
      // 即使API调用失败，也在本地保存，确保用户体验
      const lastMessage = messages.value[messages.value.length - 1]
      const memory = {
        id: Date.now(),
        preview: lastMessage.content.length > 30 ? lastMessage.content.substring(0, 30) + '...' : lastMessage.content,
        timestamp: new Date().toLocaleString('zh-CN'),
        messages: [...messages.value]
      }
      conversationMemories.value.unshift(memory)
      
      // 限制记忆数量
      if (conversationMemories.value.length > 10) {
        conversationMemories.value.pop()
      }
    }
  }
}

const loadMemory = async (memory) => {
  try {
    // 从后端API加载特定对话记忆
    const response = await axios.get(`http://localhost:8000/api/v1/chats/${memory.id}`)
    if (response.data.data) {
      currentChatId.value = response.data.data.id;
      // 加载消息列表
      const messagesResponse = await axios.get(`http://localhost:8000/api/v1/chats/${memory.id}/messages`)
      if (messagesResponse.data.data) {
        // 映射消息并添加创建时间戳用于排序
        const mappedMessages = messagesResponse.data.data.map(msg => ({
          id: msg.id,
          role: msg.sender_type === 'user' ? 'user' : 'bot',
          content: msg.content,
          timestamp: new Date(msg.created_at).toLocaleString('zh-CN'),
          createdAt: new Date(msg.created_at).getTime() // 添加时间戳用于排序
        }));
        
        // 按照创建时间从早到晚排序（升序）
        mappedMessages.sort((a, b) => a.createdAt - b.createdAt);
        
        // 移除createdAt字段并设置到messages数组
        messages.value = mappedMessages.map(({ createdAt, ...rest }) => rest);
      }
    }
  } catch (error) {
    // 降级方案：使用传入的memory对象
    if (memory.messages && memory.messages.length > 0) {
      // 如果本地消息中有序号，可以按序号排序
      messages.value = [...memory.messages].sort((a, b) => {
        // 如果有id字段，可以尝试按id排序
        if (a.id && b.id) {
          return Number(a.id) - Number(b.id);
        }
        return 0;
      });
    } else {
      messages.value = [];
    }
  }
  // 切换到新的对话时滚动到底部
  scrollToBottom()
}

const scrollToBottom = async () => {
  await nextTick()
  if (chatHistory.value) {
    chatHistory.value.scrollTop = chatHistory.value.scrollHeight
  }
}

const goBack = () => {
  router.push({ name: 'StarSelection' })
}

const startNewChat = () => {
  // 保存当前对话到记忆（如果有消息）
  if (messages.value.length > 0) {
    saveConversationToMemory()
  }
  // 清空当前对话消息和会话ID
  messages.value = []
  currentChatId.value = null
  // 切换到新的对话时滚动到底部
  scrollToBottom()
}



const toggleSidebar = () => {
  sidebarVisible.value = !sidebarVisible.value
}

// 添加侧边栏显示/隐藏的类名处理
const sidebarClass = ref('sidebar')
watch(sidebarVisible, (newVal) => {
  sidebarClass.value = newVal ? 'sidebar active' : 'sidebar'
})

watch(() => route.params.starId, (newStarId) => {
  starId.value = newStarId
  messages.value = []
  loadStarInfo()
  loadConversationMemories()
})
</script>

<style scoped>
/* 全局样式重置 */
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.chat-page {
  display: flex;
  height: 100vh;
  background: #f5f7fa;
  position: relative;
}

/* 侧边栏样式 */
.sidebar {
  width: 300px;
  background: #fff;
  border-right: 1px solid #e0e0e0;
  display: flex;
  flex-direction: column;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.05);
  z-index: 20;
  transition: transform 0.3s ease;
}

.sidebar-header {
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.back-button-mobile {
  display: none;
  background: none;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  color: #333;
}

.sidebar-header h2 {
  font-size: 1.2rem;
  font-weight: 600;
  color: #333;
}

/* 明星信息样式 */
.star-info {
  padding: 1.5rem;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.star-avatar-large {
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  font-weight: 600;
}

.star-details h3 {
  font-size: 1.1rem;
  margin-bottom: 0.25rem;
  color: #333;
}

.star-details p {
  font-size: 0.85rem;
  color: #666;
  line-height: 1.4;
}

/* 对话记忆样式 */
.conversation-memories {
  flex: 1;
  padding: 1rem;
  overflow-y: auto;
}

.conversation-memories h3 {
  font-size: 0.9rem;
  font-weight: 600;
  color: #666;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 1rem;
  padding: 0 0.5rem;
}

.memory-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.memory-item {
  background: #f8f9fa;
  padding: 1rem;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 2px solid transparent;
}

.memory-item:hover {
  background: #f0f2f5;
  border-color: #007bff;
  transform: translateX(4px);
}

.memory-preview {
  font-size: 0.9rem;
  color: #333;
  margin-bottom: 0.5rem;
  line-height: 1.3;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.memory-time {
  font-size: 0.75rem;
  color: #999;
}

/* 主聊天区域样式 */
.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  position: relative;
  width: 800px;
  overflow: hidden;
}

.chat-header {
  background: #fff;
  padding: 1rem 2rem;
  border-bottom: 1px solid #e9ecef;
  display: flex;
  align-items: center;
  gap: 1rem;
  position: sticky;
  top: 0;
  z-index: 10;
}

.back-button {
  background: none;
  border: none;
  font-size: 1.2rem;
  cursor: pointer;
  color: #333;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  transition: background-color 0.2s ease;
}

.back-button:hover {
  background-color: #f0f0f0;
}

.header-info {
  flex: 1;
}

.header-info h2 {
  font-size: 1.1rem;
  font-weight: 600;
  color: #333;
}

.header-actions {
  display: flex;
  gap: 0.5rem;
}

.new-chat-button {
  background: #f0f0f0;
  border: 1px solid #e0e0e0;
  color: #333;
  padding: 0.6rem 1.2rem;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.3s ease;
}

.new-chat-button:hover {
  background: #e0e0e0;
  border-color: #d0d0d0;
}

/* 聊天历史容器 */
.chat-history-container {
  flex: 1;
  background: #ffffff;
  padding: 1rem 0;
  overflow: hidden;
}

/* 聊天历史样式 */
.chat-history {
  width: 100%;
  height: 100%;
  overflow-y: auto;
  padding: 1rem 2rem;
  background-color: #ffffff;
}

/* 欢迎消息样式 */
.welcome-message {
  width: 100%;
  text-align: center;
  color: #6c757d;
  font-size: 1.1rem;
  padding: 3rem 0;
  min-height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.welcome-avatar {
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
  font-weight: 600;
  margin: 0 auto 1.5rem;
}

.welcome-message h3 {
  font-size: 1.5rem;
  color: #333;
  margin-bottom: 0.5rem;
}

.welcome-message p {
  font-size: 1rem;
  margin-bottom: 2rem;
  color: #666;
}

.welcome-tip {
  background: #f8f9fa;
  padding: 1rem;
  border-radius: 8px;
  border: 1px solid #e9ecef;
  display: inline-block;
}

.welcome-tip p {
  margin: 0;
  font-size: 0.9rem;
  color: #666;
}

/* 消息样式 */
.message {
  margin-bottom: 1.5rem;
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
}

.message.bot {
  justify-content: flex-start;
}

.message.user {
  justify-content: flex-end;
}

.message-avatar {
  width: 40px;
  height: 40px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9rem;
  font-weight: 600;
  flex-shrink: 0;
  margin-top: 0.25rem;
}

.message-avatar.user-avatar {
  background: #007bff;
  order: 1;
}

.message-wrapper {
  max-width: 80%;
  position: relative;
}

.message.bot .message-wrapper {
  order: 1;
}

.message.user .message-wrapper {
  order: 0;
}

.message-content {
  padding: 0.875rem 1.125rem;
  border-radius: 12px;
  word-wrap: break-word;
  line-height: 1.5;
  font-size: 1rem;
  background: #ffffff;
  border: 1px solid #e9ecef;
  text-align: left;
  /* 实现第一行最长第二行较短的效果 */
  max-width: 90%;
  display: inline-block;
  margin-bottom: 0.25rem;
}

.message.user .message-content {
  background: #007bff;
  color: white;
  border: 1px solid #007bff;
}

.message-time {
  font-size: 0.75rem;
  color: #999;
  margin-top: 0.25rem;
  text-align: left;
  margin-left: 0;
}

.message.user .message-time {
  text-align: right;
}

/* 加载状态样式 */
.message.bot.loading {
  justify-content: flex-start;
}

.loading-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #666;
  margin: 0 2px;
  animation: loading 1.4s infinite ease-in-out both;
}

.loading-dot:nth-child(1) {
  animation-delay: -0.32s;
}

.loading-dot:nth-child(2) {
  animation-delay: -0.16s;
}

@keyframes loading {
  0%, 80%, 100% {
    transform: scale(0);
  }
  40% {
    transform: scale(1);
  }
}

/* 输入区域样式 */
.chat-input-container {
  padding: 0.75rem 2rem;
  background: #fff;
  border-top: 1px solid #e9ecef;
  display: flex;
  justify-content: center;
}

.input-wrapper {
  flex: 1;
  max-width: 800px;
  position: relative;
  border: 1px solid #e0e0e0;
  border-radius: 12px;
  background: #ffffff;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.chat-input {
  width: 100%;
  min-height: 48px;
  max-height: 200px;
  padding: 0.75rem 1rem;
  border: none;
  outline: none;
  font-size: 1rem;
  background: transparent;
  resize: none;
  font-family: inherit;
  line-height: 1.5;
}

.chat-input:focus {
  background: transparent;
}

/* 输入框按钮区域 */
.input-actions {
  display: flex;
  gap: 0.5rem;
  padding: 0 1rem 0.75rem 1rem;
  border-top: 1px solid #f0f0f0;
}

.input-right-actions {
  position: absolute;
  right: 0.5rem;
  bottom: 0.5rem;
  display: flex;
  gap: 0.5rem;
}

.action-button {
  padding: 0.375rem 0.75rem;
  background: #f8f9fa;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.875rem;
  color: #6c757d;
  display: flex;
  align-items: center;
  gap: 0.25rem;
  transition: all 0.2s ease;
}

.action-button:hover {
  background: #e9ecef;
}

.send-button-icon {
  width: 36px;
  height: 36px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.send-button-icon:hover:not(:disabled) {
  background: #0056b3;
  transform: scale(1.05);
}

.send-button-icon:disabled {
  background: #c0c0c0;
  cursor: not-allowed;
}

/* 图标样式 */
.icon {
  font-size: 0.875rem;
}

.send-button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

/* 移动端侧边栏切换按钮 */
.sidebar-toggle {
  display: none;
  position: fixed;
  bottom: 1rem;
  right: 1rem;
  width: 56px;
  height: 56px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 50%;
  font-size: 1.2rem;
  cursor: pointer;
  box-shadow: 0 2px 10px rgba(0, 123, 255, 0.3);
  z-index: 30;
  transition: all 0.3s ease;
}

.sidebar-toggle:hover {
  background: #0056b3;
  transform: scale(1.05);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    left: 0;
    top: 0;
    height: 100vh;
    transform: translateX(-100%);
  }
  
  .sidebar.active {
    transform: translateX(0);
  }
  
  .back-button-mobile {
    display: block;
  }
  
  .sidebar-toggle {
    display: block;
  }
  
  .chat-header .back-button {
    display: none;
  }
  
  .message-wrapper {
    max-width: 85%;
  }
  
  .chat-history {
    padding: 1rem;
  }
}
</style>