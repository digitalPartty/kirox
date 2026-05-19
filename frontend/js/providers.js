// ===== DuckDuckGo / TEmail / DirectMail 配置管理 =====

let duckConfigs = [];
let temailConfigs = [];
let directMailConfigs = [];

// ===== DuckDuckGo =====

async function loadDuckConfigs() {
  try {
    duckConfigs = await window.go.main.App.GetDuckDuckGoConfigs() || [];
  } catch(e) { duckConfigs = []; }
  renderDuckList();
  updateDuckSummary();
}

function updateDuckSummary() {
  const el = document.getElementById('settings-duck-summary');
  if (el) el.textContent = duckConfigs.length > 0 ? '已配置 ' + duckConfigs.length + ' 个' : '未配置';
}

function renderDuckList() {
  const list = document.getElementById('duck-inline-list');
  if (!list) return;
  if (duckConfigs.length === 0) {
    list.innerHTML = '<div class="moemail-empty-state"><div>暂无 Token</div></div>';
    return;
  }
  list.innerHTML = duckConfigs.map(function(cfg, idx) {
    var masked = cfg.token.length > 12 ? cfg.token.slice(0, 6) + '...' + cfg.token.slice(-4) : cfg.token;
    return '<div class="moemail-config-item"><div class="moemail-config-main"><div class="moemail-status-dot success"></div><div class="moemail-config-info"><div class="moemail-config-name">Token #' + (idx+1) + '</div><div class="moemail-config-details"><span class="moemail-config-url">' + escapeHtml(masked) + '</span></div></div></div><div class="moemail-config-actions"><button onclick="deleteDuckConfig(' + idx + ')" class="btn btn-secondary btn-sm" style="color:var(--danger);">删除</button></div></div>';
  }).join('');
}

async function inlineTestDuck() {
  var token = document.getElementById('duck-inline-token').value.trim();
  if (!token) { showToast('请输入 Token', 'error'); return; }
  var btn = document.getElementById('duck-inline-test-btn');
  var st = document.getElementById('duck-inline-status');
  btn.disabled = true; btn.textContent = '测试中...';
  st.textContent = '';
  try {
    var r = await window.go.main.App.TestDuckDuckGoConnection(token);
    if (r.success) { st.style.color = 'var(--success)'; st.textContent = '成功，别名: ' + r.alias; }
    else { st.style.color = 'var(--danger)'; st.textContent = r.error || '失败'; }
  } catch(e) { st.style.color = 'var(--danger)'; st.textContent = '测试失败'; }
  btn.disabled = false; btn.textContent = '测试';
}

async function inlineAddDuck() {
  var token = document.getElementById('duck-inline-token').value.trim();
  if (!token) { showToast('请输入 Token', 'error'); return; }
  if (duckConfigs.some(function(c) { return c.token === token; })) { showToast('Token 已存在', 'error'); return; }
  duckConfigs.push({token: token});
  var r = await window.go.main.App.SaveDuckDuckGoConfigs(JSON.stringify(duckConfigs));
  if (r.error) { duckConfigs.pop(); showToast('保存失败: ' + r.error, 'error'); return; }
  document.getElementById('duck-inline-token').value = '';
  showToast('已添加');
  renderDuckList(); updateDuckSummary();
}

async function deleteDuckConfig(idx) {
  duckConfigs.splice(idx, 1);
  await window.go.main.App.SaveDuckDuckGoConfigs(JSON.stringify(duckConfigs));
  renderDuckList(); updateDuckSummary();
  showToast('已删除');
}

// ===== TEmail =====

async function loadTEmailConfigs() {
  try {
    temailConfigs = await window.go.main.App.GetTEmailConfigs() || [];
  } catch(e) { temailConfigs = []; }
  renderTEmailList();
  updateTEmailSummary();
}

function updateTEmailSummary() {
  const el = document.getElementById('settings-temail-summary');
  if (el) el.textContent = temailConfigs.length > 0 ? '已配置 ' + temailConfigs.length + ' 个' : '未配置';
}

function renderTEmailList() {
  const list = document.getElementById('temail-inline-list');
  if (!list) return;
  if (temailConfigs.length === 0) {
    list.innerHTML = '<div class="moemail-empty-state"><div>暂无配置</div></div>';
    return;
  }
  list.innerHTML = temailConfigs.map(function(cfg, idx) {
    return '<div class="moemail-config-item"><div class="moemail-config-main"><div class="moemail-status-dot success"></div><div class="moemail-config-info"><div class="moemail-config-name">' + escapeHtml(cfg.name) + '</div><div class="moemail-config-details"><span class="moemail-config-url">' + escapeHtml(cfg.baseUrl) + ' | ' + escapeHtml(cfg.email) + '</span></div></div></div><div class="moemail-config-actions"><button onclick="testTEmailByIdx(' + idx + ')" class="btn btn-secondary btn-sm">测试</button><button onclick="deleteTEmailConfig(' + idx + ')" class="btn btn-secondary btn-sm" style="color:var(--danger);">删除</button></div></div>';
  }).join('');
}

async function inlineTestTEmail() {
  var url = document.getElementById('temail-inline-url').value.trim();
  var email = document.getElementById('temail-inline-email').value.trim();
  var jwt = document.getElementById('temail-inline-jwt').value.trim();
  if (!url || !email || !jwt) { showToast('请填写完整', 'error'); return; }
  var btn = document.getElementById('temail-inline-test-btn');
  var st = document.getElementById('temail-inline-status');
  btn.disabled = true; btn.textContent = '测试中...'; st.textContent = '';
  var cfg = {name:'test', baseUrl:url, email:email, jwt:'', adminPassword:''};
  if (jwt.length > 50) { cfg.jwt = jwt; } else { cfg.adminPassword = jwt; }
  try {
    var r = await window.go.main.App.TestTEmailConnection(JSON.stringify(cfg));
    if (r.success) { st.style.color = 'var(--success)'; st.textContent = '连接成功'; }
    else { st.style.color = 'var(--danger)'; st.textContent = r.error || '失败'; }
  } catch(e) { st.style.color = 'var(--danger)'; st.textContent = '测试失败'; }
  btn.disabled = false; btn.textContent = '测试';
}

async function inlineAddTEmail() {
  var name = document.getElementById('temail-inline-name').value.trim();
  var url = document.getElementById('temail-inline-url').value.trim();
  var email = document.getElementById('temail-inline-email').value.trim();
  var jwt = document.getElementById('temail-inline-jwt').value.trim();
  if (!url || !email || !jwt) { showToast('请填写完整', 'error'); return; }
  if (!name) name = 'TEmail ' + (temailConfigs.length + 1);
  var cfg = {name:name, baseUrl:url, email:email, jwt:'', adminPassword:''};
  if (jwt.length > 50) { cfg.jwt = jwt; } else { cfg.adminPassword = jwt; }
  temailConfigs.push(cfg);
  var r = await window.go.main.App.SaveTEmailConfigs(JSON.stringify(temailConfigs));
  if (r.error) { temailConfigs.pop(); showToast('保存失败: ' + r.error, 'error'); return; }
  document.getElementById('temail-inline-name').value = '';
  document.getElementById('temail-inline-url').value = '';
  document.getElementById('temail-inline-email').value = '';
  document.getElementById('temail-inline-jwt').value = '';
  showToast('已添加: ' + name);
  renderTEmailList(); updateTEmailSummary();
}

async function testTEmailByIdx(idx) {
  var cfg = temailConfigs[idx];
  try {
    var r = await window.go.main.App.TestTEmailConnection(JSON.stringify(cfg));
    showToast(cfg.name + ': ' + (r.success ? '连接成功' : (r.error || '失败')), r.success ? 'success' : 'error');
  } catch(e) { showToast(cfg.name + ': 测试失败', 'error'); }
}

async function deleteTEmailConfig(idx) {
  temailConfigs.splice(idx, 1);
  await window.go.main.App.SaveTEmailConfigs(JSON.stringify(temailConfigs));
  renderTEmailList(); updateTEmailSummary();
  showToast('已删除');
}

// ===== DirectMail =====

async function loadDirectMailConfigs() {
  try {
    directMailConfigs = await window.go.main.App.GetDirectMailConfigs() || [];
  } catch(e) { directMailConfigs = []; }
  renderDirectMailList();
  updateDirectMailSummary();
}

function updateDirectMailSummary() {
  const el = document.getElementById('settings-directmail-summary');
  if (el) el.textContent = directMailConfigs.length > 0 ? '已配置 ' + directMailConfigs.length + ' 个' : '未配置';
}

function renderDirectMailList() {
  const list = document.getElementById('dm-inline-list');
  if (!list) return;
  if (directMailConfigs.length === 0) {
    list.innerHTML = '<div class="moemail-empty-state"><div>暂无配置</div></div>';
    return;
  }
  list.innerHTML = directMailConfigs.map(function(cfg, idx) {
    return '<div class="moemail-config-item"><div class="moemail-config-main"><div class="moemail-status-dot success"></div><div class="moemail-config-info"><div class="moemail-config-name">' + escapeHtml(cfg.name) + '</div><div class="moemail-config-details"><span class="moemail-config-url">' + escapeHtml(cfg.email) + ' | ' + escapeHtml(cfg.baseUrl) + '</span></div></div></div><div class="moemail-config-actions"><button onclick="testDirectMailByIdx(' + idx + ')" class="btn btn-secondary btn-sm">测试</button><button onclick="deleteDirectMailConfig(' + idx + ')" class="btn btn-secondary btn-sm" style="color:var(--danger);">删除</button></div></div>';
  }).join('');
}

async function inlineTestDirectMail() {
  var url = document.getElementById('dm-inline-url').value.trim();
  var email = document.getElementById('dm-inline-email').value.trim();
  var clientId = document.getElementById('dm-inline-clientid').value.trim();
  var refreshToken = document.getElementById('dm-inline-refreshtoken').value.trim();
  if (!url || !email || !clientId || !refreshToken) { showToast('请填写完整', 'error'); return; }
  var btn = document.getElementById('dm-inline-test-btn');
  var st = document.getElementById('dm-inline-status');
  btn.disabled = true; btn.textContent = '测试中...'; st.textContent = '';
  var cfg = {name:'test', baseUrl:url, email:email, clientId:clientId, refreshToken:refreshToken, mailbox:'INBOX'};
  try {
    var r = await window.go.main.App.TestDirectMailConnection(JSON.stringify(cfg));
    if (r.success) { st.style.color = 'var(--success)'; st.textContent = '连接成功'; }
    else { st.style.color = 'var(--danger)'; st.textContent = r.error || '失败'; }
  } catch(e) { st.style.color = 'var(--danger)'; st.textContent = '测试失败'; }
  btn.disabled = false; btn.textContent = '测试';
}

async function inlineAddDirectMail() {
  var name = document.getElementById('dm-inline-name').value.trim();
  var url = document.getElementById('dm-inline-url').value.trim();
  var email = document.getElementById('dm-inline-email').value.trim();
  var clientId = document.getElementById('dm-inline-clientid').value.trim();
  var refreshToken = document.getElementById('dm-inline-refreshtoken').value.trim();
  if (!url || !email || !clientId || !refreshToken) { showToast('请填写完整', 'error'); return; }
  if (!name) name = 'DirectMail ' + (directMailConfigs.length + 1);
  var cfg = {name:name, baseUrl:url, email:email, clientId:clientId, refreshToken:refreshToken, mailbox:'INBOX'};
  directMailConfigs.push(cfg);
  var r = await window.go.main.App.SaveDirectMailConfigs(JSON.stringify(directMailConfigs));
  if (r.error) { directMailConfigs.pop(); showToast('保存失败: ' + r.error, 'error'); return; }
  document.getElementById('dm-inline-name').value = '';
  document.getElementById('dm-inline-url').value = '';
  document.getElementById('dm-inline-email').value = '';
  document.getElementById('dm-inline-clientid').value = '';
  document.getElementById('dm-inline-refreshtoken').value = '';
  showToast('已添加: ' + name);
  renderDirectMailList(); updateDirectMailSummary();
}

async function testDirectMailByIdx(idx) {
  var cfg = directMailConfigs[idx];
  try {
    var r = await window.go.main.App.TestDirectMailConnection(JSON.stringify(cfg));
    showToast(cfg.name + ': ' + (r.success ? '连接成功' : (r.error || '失败')), r.success ? 'success' : 'error');
  } catch(e) { showToast(cfg.name + ': 测试失败', 'error'); }
}

async function deleteDirectMailConfig(idx) {
  directMailConfigs.splice(idx, 1);
  await window.go.main.App.SaveDirectMailConfigs(JSON.stringify(directMailConfigs));
  renderDirectMailList(); updateDirectMailSummary();
  showToast('已删除');
}

// ===== 初始化 =====
document.addEventListener('DOMContentLoaded', function() {
  loadDuckConfigs();
  loadTEmailConfigs();
  loadDirectMailConfigs();
});
