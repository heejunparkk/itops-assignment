<template>
  <div id="app">
    <header class="header">
      <div class="container">
        <h1>이슈 관리 시스템</h1>
      </div>
    </header>
    <main class="main">
      <div class="container">
        <router-view />
      </div>
    </main>
  </div>
</template>

<script>
export default {
  name: 'App',
}
</script>

<style>
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: 'Noto Sans KR', Arial, sans-serif;
  line-height: 1.6;
  color: #333;
  background-color: #f5f5f5;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 15px;
}

.header {
  background-color: #4a6baf;
  color: white;
  padding: 1rem 0;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header h1 {
  margin: 0;
  font-size: 1.5rem;
}

.main {
  padding: 2rem 0;
}

.btn {
  display: inline-block;
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  font-size: 0.9rem;
  cursor: pointer;
  text-decoration: none;
  transition: background-color 0.3s;
}

.btn-primary {
  background-color: #4a6baf;
  color: white;
}

.btn-primary:hover {
  background-color: #3a5a9f;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background-color: #5a6268;
}

.card {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
  margin-bottom: 1.5rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
}

.form-control {
  width: 100%;
  padding: 0.5rem;
  font-size: 1rem;
  border: 1px solid #ced4da;
  border-radius: 4px;
}

.form-control:focus {
  border-color: #4a6baf;
  outline: none;
}

.status-badge {
  display: inline-block;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.8rem;
  font-weight: 500;
}

.status-pending {
  background-color: #ffc107;
  color: #212529;
}

.status-in-progress {
  background-color: #17a2b8;
  color: white;
}

.status-completed {
  background-color: #28a745;
  color: white;
}

.status-cancelled {
  background-color: #dc3545;
  color: white;
}
</style>
