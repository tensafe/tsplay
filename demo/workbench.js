(function () {
  const state = {
    artifactRoot: "",
    appMeta: null,
    sessions: [],
    providers: [],
    sites: [],
    currentSiteId: "",
    currentTab: "pages",
    pages: [],
    apis: [],
    entities: [],
    lastExplore: null,
    lastPlan: null,
    lastRun: null,
    lastRepair: null,
  };

  const refs = {};

  document.addEventListener("DOMContentLoaded", () => {
    cacheRefs();
    bindEvents();
    applyProviderTypeDefaults();
    runAction("正在加载 Workbench", refreshAll);
  });

  function cacheRefs() {
    refs.artifactRootLabel = document.getElementById("artifactRootLabel");
    refs.currentSiteLabel = document.getElementById("currentSiteLabel");
    refs.statusBar = document.getElementById("statusBar");
    refs.refreshAllButton = document.getElementById("refreshAllButton");
    refs.sitePicker = document.getElementById("sitePicker");
    refs.loadSiteButton = document.getElementById("loadSiteButton");
    refs.siteForm = document.getElementById("siteForm");
    refs.siteIdInput = document.getElementById("siteIdInput");
    refs.siteNameInput = document.getElementById("siteNameInput");
    refs.siteUrlInput = document.getElementById("siteUrlInput");
    refs.siteDomainsInput = document.getElementById("siteDomainsInput");
    refs.siteSessionSelect = document.getElementById("siteSessionSelect");
    refs.sessionSummary = document.getElementById("sessionSummary");
    refs.siteProviderSelect = document.getElementById("siteProviderSelect");
    refs.providerSummary = document.getElementById("providerSummary");
    refs.newSiteButton = document.getElementById("newSiteButton");
    refs.sessionForm = document.getElementById("sessionForm");
    refs.sessionNameInput = document.getElementById("sessionNameInput");
    refs.sessionStoragePathInput = document.getElementById("sessionStoragePathInput");
    refs.sessionProfileInput = document.getElementById("sessionProfileInput");
    refs.sessionProfileSessionInput = document.getElementById("sessionProfileSessionInput");
    refs.providerForm = document.getElementById("providerForm");
    refs.providerIdInput = document.getElementById("providerIdInput");
    refs.providerNameInput = document.getElementById("providerNameInput");
    refs.providerTypeSelect = document.getElementById("providerTypeSelect");
    refs.providerBaseURLInput = document.getElementById("providerBaseURLInput");
    refs.providerModelInput = document.getElementById("providerModelInput");
    refs.providerAPIKeyEnvInput = document.getElementById("providerAPIKeyEnvInput");
    refs.providerAPIKeyInput = document.getElementById("providerAPIKeyInput");
    refs.providerSystemPromptInput = document.getElementById("providerSystemPromptInput");
    refs.providerEnabledInput = document.getElementById("providerEnabledInput");
    refs.exploreForm = document.getElementById("exploreForm");
    refs.exploreMaxPages = document.getElementById("exploreMaxPages");
    refs.exploreTimeout = document.getElementById("exploreTimeout");
    refs.exploreHeadless = document.getElementById("exploreHeadless");
    refs.exploreMeta = document.getElementById("exploreMeta");
    refs.cardList = document.getElementById("cardList");
    refs.tabs = Array.from(document.querySelectorAll(".wb-tab"));
    refs.taskForm = document.getElementById("taskForm");
    refs.intentInput = document.getElementById("intentInput");
    refs.runHeadlessInput = document.getElementById("runHeadlessInput");
    refs.runFlowButton = document.getElementById("runFlowButton");
    refs.planSummary = document.getElementById("planSummary");
    refs.planCandidates = document.getElementById("planCandidates");
    refs.flowEditor = document.getElementById("flowEditor");
    refs.copyFlowButton = document.getElementById("copyFlowButton");
    refs.runResultPanel = document.getElementById("runResultPanel");
    refs.replayFlowButton = document.getElementById("replayFlowButton");
    refs.repairFlowButton = document.getElementById("repairFlowButton");
    refs.autoRepairFlowButton = document.getElementById("autoRepairFlowButton");
    refs.runRepairedFlowButton = document.getElementById("runRepairedFlowButton");
    refs.repairPanel = document.getElementById("repairPanel");
    refs.runTraceList = document.getElementById("runTraceList");
  }

  function bindEvents() {
    refs.refreshAllButton.addEventListener("click", () => runAction("正在刷新数据", refreshAll));
    refs.loadSiteButton.addEventListener("click", () => runAction("正在加载站点", () => selectSite(refs.sitePicker.value)));
    refs.sitePicker.addEventListener("change", () => runAction("正在切换站点", () => selectSite(refs.sitePicker.value)));
    refs.siteSessionSelect.addEventListener("change", renderSessionSummary);
    refs.siteProviderSelect.addEventListener("change", renderProviderSummary);
    refs.newSiteButton.addEventListener("click", resetSiteForm);
    refs.siteForm.addEventListener("submit", (event) => {
      event.preventDefault();
      runAction("正在保存站点", saveSite);
    });
    refs.sessionForm.addEventListener("submit", (event) => {
      event.preventDefault();
      runAction("正在保存 Session", saveSession);
    });
    refs.providerForm.addEventListener("submit", (event) => {
      event.preventDefault();
      runAction("正在保存 Provider", saveProvider);
    });
    refs.providerTypeSelect.addEventListener("change", applyProviderTypeDefaults);
    refs.exploreForm.addEventListener("submit", (event) => {
      event.preventDefault();
      runAction("正在探索站点", exploreSite);
    });
    refs.taskForm.addEventListener("submit", (event) => {
      event.preventDefault();
      runAction("正在生成 Flow", planTask);
    });
    refs.runFlowButton.addEventListener("click", () => runAction("正在执行 Flow", executePlanFlow));
    refs.replayFlowButton.addEventListener("click", () => runAction("正在回放 Flow", executePlanFlow));
    refs.repairFlowButton.addEventListener("click", () => runAction("正在生成 Repair Context", buildRepairContext));
    refs.autoRepairFlowButton.addEventListener("click", () => runAction("正在自动修复 Flow", buildAutoRepair));
    refs.runRepairedFlowButton.addEventListener("click", () => runAction("正在执行修复版 Flow", executePlanFlow));
    refs.copyFlowButton.addEventListener("click", () => runAction("正在复制 Flow", copyFlow));
    refs.tabs.forEach((button) => {
      button.addEventListener("click", () => {
        state.currentTab = button.dataset.tab || "pages";
        renderTabs();
        renderCards();
      });
    });
  }

  async function runAction(label, action) {
    setStatus(label + "…", "busy");
    try {
      await action();
      setStatus(label + "完成", "idle");
    } catch (error) {
      console.error(error);
      setStatus(error.message || String(error), "error");
    }
  }

  async function refreshAll() {
    const [health, appMeta] = await Promise.all([
      apiFetch("/api/workbench/health"),
      apiFetchOptional("/api/workbench/app-meta", null),
    ]);
    state.artifactRoot = health.artifact_root || "";
    state.appMeta = appMeta;
    refs.artifactRootLabel.textContent =
      (state.appMeta && state.appMeta.artifact_root) || state.artifactRoot || "default";

    await Promise.all([loadSessions(), loadProviders(), loadSites()]);
    if (state.currentSiteId) {
      await loadKnowledge(state.currentSiteId);
    } else {
      renderExploreMeta();
      renderCards();
    }
    renderPlan();
    renderRunResult();
    renderRepair();
  }

  async function loadSessions(selectedName) {
    const payload = await apiFetch("/api/workbench/sessions");
    state.sessions = Array.isArray(payload.sessions) ? payload.sessions : [];
    if (!selectedName && refs.siteSessionSelect.value) {
      selectedName = refs.siteSessionSelect.value;
    }
    const names = state.sessions.map((item) => item.name);
    if (!selectedName || !names.includes(selectedName)) {
      selectedName = names[0] || "";
    }
    refs.siteSessionSelect.innerHTML = buildOptions(
      [{ value: "", label: "不使用已保存 Session" }].concat(
        state.sessions.map((item) => ({
          value: item.name,
          label: item.name + " · " + (item.kind || "session"),
        }))
      ),
      selectedName
    );
    renderSessionSummary();
  }

  async function loadProviders(selectedProviderId) {
    const payload = await apiFetch("/api/workbench/providers");
    state.providers = Array.isArray(payload.providers) ? payload.providers : [];
    if (!selectedProviderId && refs.siteProviderSelect.value) {
      selectedProviderId = refs.siteProviderSelect.value;
    }
    const providerIds = state.providers.map((item) => item.provider_id);
    if (!selectedProviderId || !providerIds.includes(selectedProviderId)) {
      selectedProviderId = "";
    }
    refs.siteProviderSelect.innerHTML = buildOptions(
      [{ value: "", label: "自动选择已就绪 Provider" }].concat(
        state.providers.map((item) => ({
          value: item.provider_id,
          label: describeProviderOption(item),
        }))
      ),
      selectedProviderId
    );
    renderProviderSummary();
  }

  async function loadSites(selectedSiteId) {
    const payload = await apiFetch("/api/workbench/sites");
    state.sites = Array.isArray(payload.sites) ? payload.sites : [];
    const siteIds = state.sites.map((item) => item.site_id);
    const nextSiteId =
      normalizeString(selectedSiteId) ||
      normalizeString(state.currentSiteId) ||
      siteIds[0] ||
      "";

    refs.sitePicker.innerHTML = buildOptions(
      [{ value: "", label: "请选择站点" }].concat(
        state.sites.map((item) => ({
          value: item.site_id,
          label: (item.name || item.site_id) + " · " + item.site_id,
        }))
      ),
      nextSiteId
    );

    state.currentSiteId = siteIds.includes(nextSiteId) ? nextSiteId : "";
    if (state.currentSiteId) {
      populateSiteForm(findSite(state.currentSiteId));
    } else {
      refs.currentSiteLabel.textContent = "未选择";
    }
  }

  async function selectSite(siteId) {
    siteId = normalizeString(siteId);
    state.currentSiteId = siteId;
    refs.sitePicker.value = siteId;
    const site = findSite(siteId);
    populateSiteForm(site);
    state.lastPlan = null;
    state.lastRun = null;
    state.lastRepair = null;
    renderPlan();
    renderRunResult();
    renderRepair();
    if (!siteId) {
      state.pages = [];
      state.apis = [];
      state.entities = [];
      state.lastExplore = null;
      renderExploreMeta();
      renderCards();
      refs.currentSiteLabel.textContent = "未选择";
      return;
    }
    await loadKnowledge(siteId);
  }

  async function loadKnowledge(siteId) {
    const [pagesPayload, apisPayload, entitiesPayload] = await Promise.all([
      apiFetchOptional("/api/workbench/sites/" + encodeURIComponent(siteId) + "/pages", { pages: [] }),
      apiFetchOptional("/api/workbench/sites/" + encodeURIComponent(siteId) + "/apis", { apis: [] }),
      apiFetchOptional("/api/workbench/sites/" + encodeURIComponent(siteId) + "/entities", { entities: [] }),
    ]);
    state.pages = Array.isArray(pagesPayload.pages) ? pagesPayload.pages : [];
    state.apis = Array.isArray(apisPayload.apis) ? apisPayload.apis : [];
    state.entities = Array.isArray(entitiesPayload.entities) ? entitiesPayload.entities : [];
    state.lastExplore = {
      run_id: state.lastExplore ? state.lastExplore.run_id : "",
      pages: state.pages,
      apis: state.apis,
      entities: state.entities,
      finished_at: state.lastExplore ? state.lastExplore.finished_at : "",
    };
    refs.currentSiteLabel.textContent = describeSite(findSite(siteId));
    renderExploreMeta();
    renderCards();
  }

  async function saveSite() {
    const payload = {
      site_id: normalizeString(refs.siteIdInput.value),
      name: normalizeString(refs.siteNameInput.value),
      start_url: normalizeString(refs.siteUrlInput.value),
      allowed_domains: parseList(refs.siteDomainsInput.value),
      session_name: normalizeString(refs.siteSessionSelect.value),
      provider_id: normalizeString(refs.siteProviderSelect.value),
    };
    const saved = await apiFetch("/api/workbench/sites", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    await loadSites(saved.site_id);
    refs.sitePicker.value = saved.site_id;
    await selectSite(saved.site_id);
  }

  async function saveSession() {
    const payload = {
      name: normalizeString(refs.sessionNameInput.value),
      storage_state_path: normalizeString(refs.sessionStoragePathInput.value),
      profile: normalizeString(refs.sessionProfileInput.value),
      session: normalizeString(refs.sessionProfileSessionInput.value),
    };
    const saved = await apiFetch("/api/workbench/sessions", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    refs.sessionNameInput.value = "";
    refs.sessionStoragePathInput.value = "";
    refs.sessionProfileInput.value = "";
    refs.sessionProfileSessionInput.value = "";
    await loadSessions(saved.name || payload.name);
  }

  async function saveProvider() {
    const payload = {
      provider_id: normalizeString(refs.providerIdInput.value),
      name: normalizeString(refs.providerNameInput.value),
      type: normalizeString(refs.providerTypeSelect.value),
      base_url: normalizeString(refs.providerBaseURLInput.value),
      model: normalizeString(refs.providerModelInput.value),
      api_key_env: normalizeString(refs.providerAPIKeyEnvInput.value),
      api_key: normalizeString(refs.providerAPIKeyInput.value),
      system_prompt: normalizeString(refs.providerSystemPromptInput.value),
      enabled: refs.providerEnabledInput.checked,
    };
    const saved = await apiFetch("/api/workbench/providers", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    refs.providerAPIKeyInput.value = "";
    await loadProviders(saved.provider_id || payload.provider_id);
    refs.siteProviderSelect.value = saved.provider_id || payload.provider_id;
    renderProviderSummary();
  }

  async function exploreSite() {
    if (!state.currentSiteId) {
      throw new Error("请先选择一个站点");
    }
    const result = await apiFetch(
      "/api/workbench/sites/" + encodeURIComponent(state.currentSiteId) + "/explore",
      {
        method: "POST",
        body: JSON.stringify({
          headless: refs.exploreHeadless.checked,
          timeout_ms: toNumber(refs.exploreTimeout.value, 30000),
          max_pages: toNumber(refs.exploreMaxPages.value, 8),
        }),
      }
    );
    state.lastExplore = result;
    state.pages = Array.isArray(result.pages) ? result.pages : [];
    state.apis = Array.isArray(result.apis) ? result.apis : [];
    state.entities = Array.isArray(result.entities) ? result.entities : [];
    state.lastRepair = null;
    renderExploreMeta();
    renderCards();
    renderRepair();
  }

  async function planTask() {
    if (!state.currentSiteId) {
      throw new Error("请先选择一个站点");
    }
    const intent = normalizeString(refs.intentInput.value);
    if (!intent) {
      throw new Error("请输入任务需求");
    }
    const plan = await apiFetch("/api/workbench/tasks/plan", {
      method: "POST",
      body: JSON.stringify({
        site_id: state.currentSiteId,
        intent: intent,
      }),
    });
    state.lastPlan = plan;
    state.lastRun = null;
    state.lastRepair = null;
    renderPlan();
    renderRunResult();
    renderRepair();
  }

  async function executePlanFlow() {
    if (!state.currentSiteId) {
      throw new Error("请先选择一个站点");
    }
    const flowYAML = normalizeString(refs.flowEditor.value);
    if (!flowYAML || flowYAML === "还没有生成 Flow") {
      throw new Error("请先生成一个可执行的 Flow");
    }
    const payload = await apiFetch("/api/workbench/tasks/run", {
      method: "POST",
      body: JSON.stringify({
        site_id: state.currentSiteId,
        intent: normalizeString(refs.intentInput.value),
        flow_yaml: flowYAML,
        headless: refs.runHeadlessInput.checked,
      }),
    });
    state.lastRun = payload;
    state.lastRepair = null;
    if (!state.lastPlan && payload.plan) {
      state.lastPlan = payload.plan;
    }
    if (payload.flow_yaml) {
      refs.flowEditor.value = payload.flow_yaml;
    }
    renderPlan();
    renderRunResult();
    renderRepair();
    if (!payload.ok) {
      throw new Error(payload.error || "Flow 执行失败");
    }
  }

  async function buildRepairContext() {
    if (!state.lastRun) {
      throw new Error("当前还没有执行结果，先运行一次 Flow");
    }
    if (state.lastRun.ok) {
      throw new Error("当前执行已成功，不需要生成 Repair Context");
    }
    const flowYAML = normalizeString(refs.flowEditor.value);
    if (!flowYAML || flowYAML === "还没有生成 Flow") {
      throw new Error("当前没有可修复的 Flow");
    }
    const payload = await apiFetch("/api/workbench/tasks/repair", {
      method: "POST",
      body: JSON.stringify({
        artifact_root: state.artifactRoot,
        flow_yaml: flowYAML,
        run_result: state.lastRun.result || null,
        error: state.lastRun.error || "",
      }),
    });
    state.lastRepair = payload;
    renderRepair();
  }

  async function buildAutoRepair() {
    if (!state.lastRun) {
      throw new Error("当前还没有执行结果，先运行一次 Flow");
    }
    if (state.lastRun.ok) {
      throw new Error("当前执行已成功，不需要自动修复");
    }
    const flowYAML = normalizeString(refs.flowEditor.value);
    if (!flowYAML || flowYAML === "还没有生成 Flow") {
      throw new Error("当前没有可修复的 Flow");
    }
    const payload = await apiFetch("/api/workbench/tasks/repair/auto", {
      method: "POST",
      body: JSON.stringify({
        site_id: state.currentSiteId,
        provider_id: normalizeString(refs.siteProviderSelect.value),
        artifact_root: state.artifactRoot,
        flow_yaml: flowYAML,
        run_result: state.lastRun.result || null,
        error: state.lastRun.error || "",
      }),
    });
    state.lastRepair = payload;
    if (payload.repaired_flow_yaml) {
      refs.flowEditor.value = payload.repaired_flow_yaml;
    }
    renderRepair();
    if (!payload.ok) {
      throw new Error(payload.validation_error || payload.error || "自动修复失败");
    }
  }

  async function copyFlow() {
    const text = normalizeString(refs.flowEditor.value);
    if (!text) {
      throw new Error("当前没有可复制的 Flow");
    }
    await navigator.clipboard.writeText(text);
  }

  function populateSiteForm(site) {
    refs.siteIdInput.value = site && site.site_id ? site.site_id : "";
    refs.siteNameInput.value = site && site.name ? site.name : "";
    refs.siteUrlInput.value = site && site.start_url ? site.start_url : "";
    refs.siteDomainsInput.value = site && Array.isArray(site.allowed_domains) ? site.allowed_domains.join(", ") : "";
    refs.siteSessionSelect.value = site && site.session_name ? site.session_name : "";
    refs.siteProviderSelect.value = site && site.provider_id ? site.provider_id : "";
    refs.currentSiteLabel.textContent = site ? describeSite(site) : "未选择";
    renderSessionSummary();
    renderProviderSummary();
  }

  function resetSiteForm() {
    refs.siteIdInput.value = "";
    refs.siteNameInput.value = "";
    refs.siteUrlInput.value = "";
    refs.siteDomainsInput.value = "";
    refs.siteProviderSelect.value = "";
    refs.sitePicker.value = "";
    state.currentSiteId = "";
    state.pages = [];
    state.apis = [];
    state.entities = [];
    state.lastExplore = null;
    state.lastPlan = null;
    state.lastRun = null;
    state.lastRepair = null;
    renderExploreMeta();
    renderCards();
    renderPlan();
    renderRunResult();
    renderRepair();
    refs.currentSiteLabel.textContent = "未选择";
    renderProviderSummary();
  }

  function renderSessionSummary() {
    const selectedName = normalizeString(refs.siteSessionSelect.value);
    const session = state.sessions.find((item) => item.name === selectedName);
    if (!session) {
      refs.sessionSummary.textContent = "不复用登录态时，探索和 Flow 会以默认浏览器上下文执行。";
      return;
    }
    const pieces = [
      session.kind || "session",
      session.source || "",
      session.updated_at ? "更新于 " + formatDate(session.updated_at) : "",
    ].filter(Boolean);
    refs.sessionSummary.textContent = pieces.join(" · ");
  }

  function renderProviderSummary() {
    const selectedID = normalizeString(refs.siteProviderSelect.value);
    const provider = findProvider(selectedID);
    if (!provider) {
      refs.providerSummary.textContent = "未显式选择 provider 时，自动修复会优先尝试已就绪的默认 provider。";
      return;
    }
    const pieces = [
      provider.type || "provider",
      provider.resolved_model || provider.model || "",
      provider.ready ? "ready" : provider.status || "needs_config",
      provider.resolved_api_key_source || "",
      provider.detected ? "auto-detected" : "",
    ].filter(Boolean);
    refs.providerSummary.textContent =
      pieces.join(" · ") + (provider.error ? " · " + provider.error : "");
  }

  function renderTabs() {
    refs.tabs.forEach((button) => {
      button.classList.toggle("is-active", button.dataset.tab === state.currentTab);
    });
  }

  function renderExploreMeta() {
    const finishedAt = state.lastExplore && state.lastExplore.finished_at ? formatDate(state.lastExplore.finished_at) : "暂无";
    const runID = state.lastExplore && state.lastExplore.run_id ? state.lastExplore.run_id : "暂无";
    refs.exploreMeta.innerHTML = [
      metaCard("最近探索", finishedAt),
      metaCard("Run ID", runID),
      metaCard("页面数", String(state.pages.length)),
      metaCard("接口数", String(state.apis.length)),
      metaCard("实体数", String(state.entities.length)),
    ].join("");
  }

  function renderCards() {
    const renderers = {
      pages: renderPageCard,
      apis: renderAPICard,
      entities: renderEntityCard,
    };
    const collections = {
      pages: state.pages,
      apis: state.apis,
      entities: state.entities,
    };
    const items = collections[state.currentTab] || [];
    if (!items.length) {
      refs.cardList.innerHTML =
        '<article class="wb-empty-card"><h3>当前没有卡片</h3><p>先执行探索，或者换一个已经有知识沉淀的站点。</p></article>';
      return;
    }
    refs.cardList.innerHTML = items.map((item) => renderers[state.currentTab](item)).join("");
  }

  function renderPlan() {
    const plan = state.lastPlan;
    if (!plan) {
      refs.planSummary.innerHTML =
        "<h3>规划结果</h3><p>输入任务需求后，这里会展示策略、原因、候选页面/API 以及 Flow YAML。</p>";
      refs.planCandidates.innerHTML = "";
      refs.flowEditor.value = "# 还没有生成 Flow";
      refs.runFlowButton.disabled = true;
      refs.replayFlowButton.disabled = true;
      refs.autoRepairFlowButton.disabled = true;
      refs.runRepairedFlowButton.disabled = true;
      return;
    }

    const summaryBits = [
      '<div class="wb-card-meta">',
      pill("策略 · " + escapeHTML(plan.strategy || "unknown")),
      pill("Flow · " + escapeHTML(plan.flow_name || "未命名")),
      "</div>",
      "<p>" + escapeHTML(plan.reason || "没有生成解释。") + "</p>",
    ];
    refs.planSummary.innerHTML = "<h3>规划结果</h3>" + summaryBits.join("");

    const sections = [];
    if (Array.isArray(plan.matched_pages) && plan.matched_pages.length) {
      sections.push(renderCandidateSection("命中的页面", plan.matched_pages));
    }
    if (Array.isArray(plan.matched_apis) && plan.matched_apis.length) {
      sections.push(renderCandidateSection("命中的 API", plan.matched_apis));
    }
    refs.planCandidates.innerHTML = sections.join("");
    refs.flowEditor.value = plan.flow_yaml || "# 当前没有生成 Flow";
    refs.runFlowButton.disabled = !(plan.flow_yaml || "");
    refs.replayFlowButton.disabled = !(plan.flow_yaml || "");
    refs.autoRepairFlowButton.disabled = true;
    refs.runRepairedFlowButton.disabled = !(plan.flow_yaml || "");
  }

  function renderRunResult() {
    const run = state.lastRun;
    if (!run) {
      refs.runResultPanel.innerHTML =
        "<h3>执行结果</h3><p>生成 Flow 后，可以直接在这里执行，并查看 trace、变量输出和 artifact 路径。</p>";
      refs.runTraceList.innerHTML = "";
      refs.repairFlowButton.disabled = true;
      refs.autoRepairFlowButton.disabled = true;
      return;
    }

    const result = run.result || {};
    const summary = [
      '<div class="wb-card-meta">',
      pill("状态 · " + escapeHTML(run.ok ? "success" : "failed")),
      pill("Flow · " + escapeHTML(run.flow_name || result.name || "未命名")),
      (result.run_id ? pill("Run · " + escapeHTML(result.run_id)) : ""),
      (result.trace ? pill("Trace · " + escapeHTML(String(result.trace.length || 0))) : ""),
      "</div>",
      run.error ? "<p>" + escapeHTML(run.error) + "</p>" : "<p>Flow 已执行完成，可以继续查看变量输出与步骤 trace。</p>",
      renderArtifactPath("Run Root", result.run_root),
      renderArtifactPath("Browser Video", result.browser_video),
      renderSchema("Vars", result.vars),
    ].filter(Boolean);
    refs.runResultPanel.innerHTML = "<h3>执行结果</h3>" + summary.join("");
    refs.repairFlowButton.disabled = !!run.ok;
    refs.autoRepairFlowButton.disabled = !!run.ok;

    const traces = flattenTraceList(result.trace || []);
    if (!traces.length) {
      refs.runTraceList.innerHTML = "";
      return;
    }
    refs.runTraceList.innerHTML =
      '<div class="wb-trace-list">' +
      traces
        .map((trace) => {
          const artifactBits = [];
          if (trace.artifacts) {
            artifactBits.push(renderArtifactPath("Screenshot", trace.artifacts.screenshot_path));
            artifactBits.push(renderArtifactPath("HTML", trace.artifacts.html_path));
            artifactBits.push(renderArtifactPath("DOM", trace.artifacts.dom_snapshot_path));
          }
          return (
            '<article class="wb-trace">' +
            "<header><div><h4>" +
            escapeHTML((trace.path || "?") + " · " + (trace.name || trace.action || "step")) +
            "</h4><p>" +
            escapeHTML(trace.action || "") +
            "</p></div>" +
            '<div class="wb-card-meta">' +
            pill(escapeHTML(trace.status || "unknown")) +
            (trace.duration_ms != null ? pill(escapeHTML(String(trace.duration_ms)) + " ms") : "") +
            "</div></header>" +
            (trace.error ? "<p>" + escapeHTML(trace.error) + "</p>" : "") +
            (trace.output_summary ? "<p>输出：" + escapeHTML(trace.output_summary) + "</p>" : "") +
            artifactBits.join("") +
            renderSchema("Args", trace.args) +
            renderSchema("Output", trace.output) +
            "</article>"
          );
        })
        .join("") +
      "</div>";
  }

  function renderRepair() {
    const payload = state.lastRepair;
    if (!payload) {
      refs.repairPanel.innerHTML =
        "<h3>Repair 面板</h3><p>当执行失败后，这里会展示失败步骤、修复提示、相关 DOM 证据和可发送给 AI 的 repair prompt。</p>";
      return;
    }

    const context = payload.context || {};
    const repair = payload.repair || {};
    const provider = payload.provider || null;
    const hintItems = Array.isArray(context.repair_hints) ? context.repair_hints : [];
    const relevantSelectors = context.artifacts && Array.isArray(context.artifacts.relevant_selectors)
      ? context.artifacts.relevant_selectors
      : [];
    const relevantDOM = context.artifacts && Array.isArray(context.artifacts.relevant_dom)
      ? context.artifacts.relevant_dom
      : [];

    const blocks = [
      '<div class="wb-card-meta">',
      pill("失败类型 · " + escapeHTML(context.failure_category || "unknown")),
      pill("失败步骤 · " + escapeHTML(context.failed_step_path || "?")),
      (provider ? pill("Provider · " + escapeHTML(provider.name || provider.provider_id || "auto")) : ""),
      "</div>",
      context.failure_reason ? "<p>" + escapeHTML(context.failure_reason) + "</p>" : "",
      payload.validation_error
        ? "<p><strong>自动修复校验失败：</strong>" + escapeHTML(payload.validation_error) + "</p>"
        : "",
      payload.repaired_flow_yaml ? "<p>已生成修复版 Flow，并已回填到编辑器中。</p>" : "",
      renderArtifactPath("Failure Screenshot", context.artifacts && context.artifacts.paths ? context.artifacts.paths.screenshot_path : ""),
      renderArtifactPath("Failure HTML", context.artifacts && context.artifacts.paths ? context.artifacts.paths.html_path : ""),
      renderArtifactPath("Failure DOM", context.artifacts && context.artifacts.paths ? context.artifacts.paths.dom_snapshot_path : ""),
      relevantSelectors.length ? renderOptionalList("相关选择器", relevantSelectors) : "",
      relevantDOM.length ? renderOptionalList("相关 DOM 线索", relevantDOM) : "",
      hintItems.length
        ? '<div class="wb-candidate-list">' +
          hintItems
            .map((hint, index) => {
              return (
                '<section class="wb-candidate">' +
                "<header><h3>" +
                escapeHTML("Hint " + (index + 1)) +
                "</h3></header>" +
                (hint.issue ? "<p><strong>问题：</strong>" + escapeHTML(hint.issue) + "</p>" : "") +
                (hint.suggestion ? "<p><strong>建议：</strong>" + escapeHTML(hint.suggestion) + "</p>" : "") +
                (hint.selector ? "<p><strong>Selector：</strong>" + escapeHTML(hint.selector) + "</p>" : "") +
                "</section>"
              );
            })
            .join("") +
          "</div>"
        : "",
      renderSchema("Repair Context", context),
      renderSchema("Repair Request", repair),
      provider ? renderSchema("Provider", provider) : "",
      repair.prompt ? renderSchema("Repair Prompt", repair.prompt) : "",
      payload.model_output ? renderSchema("Model Output", payload.model_output) : "",
      payload.repaired_flow_yaml ? renderSchema("Repaired Flow YAML", payload.repaired_flow_yaml) : "",
    ].filter(Boolean);

    refs.repairPanel.innerHTML = "<h3>Repair 面板</h3>" + blocks.join("");
  }

  function renderCandidateSection(title, items) {
    return (
      '<section class="wb-candidate">' +
      "<header><h3>" +
      escapeHTML(title) +
      "</h3></header>" +
      items
        .map((item) => {
          const detail = item.method ? item.method + " " + (item.path || "") : item.url || "";
          return (
            '<div class="wb-card-meta">' +
            pill(escapeHTML(item.label || item.id || "candidate")) +
            pill("score " + escapeHTML(String(item.score || 0))) +
            "</div>" +
            '<p>' +
            escapeHTML(detail) +
            "</p>"
          );
        })
        .join("") +
      "</section>"
    );
  }

  function renderPageCard(page) {
    const meta = [];
    if (page.normalized_route) {
      meta.push(pill(escapeHTML(page.normalized_route)));
    }
    if (Array.isArray(page.menu_path) && page.menu_path.length) {
      meta.push(pill(escapeHTML(page.menu_path.join(" / "))));
    }
    meta.push(riskPill(page.risk));
    const actionLabels = Array.isArray(page.actions) ? page.actions.map((item) => item.label).filter(Boolean) : [];
    const tableLabels = Array.isArray(page.tables) ? page.tables.map((item) => item.name || item.selector).filter(Boolean) : [];
    const formLabels = Array.isArray(page.forms)
      ? page.forms.map((item) => {
          const fields = Array.isArray(item.fields) ? item.fields.map((field) => field.label || field.name).filter(Boolean) : [];
          return [item.name || item.selector || "表单", fields.join(" / ")].filter(Boolean).join(" · ");
        })
      : [];
    const linkLabels = Array.isArray(page.links)
      ? page.links
          .slice(0, 8)
          .map((item) => [item.text || "link", item.href || ""].filter(Boolean).join(" · "))
      : [];
    return (
      '<article class="wb-card">' +
      "<header>" +
      "<div><h3>" +
      escapeHTML(page.title || page.normalized_route || page.url || "未命名页面") +
      "</h3><p>" +
      escapeHTML(page.summary || page.url || "") +
      "</p></div>" +
      "</header>" +
      '<div class="wb-card-meta">' +
      meta.join("") +
      "</div>" +
      renderOptionalList("面包屑", page.breadcrumbs || []) +
      renderOptionalList("表单", formLabels) +
      renderOptionalList("动作", actionLabels) +
      renderOptionalList("表格", tableLabels) +
      renderOptionalList("链接", linkLabels) +
      renderArtifactPath("Observation", page.observation_path) +
      renderArtifactPath("Screenshot", page.screenshot_path) +
      renderArtifactPath("DOM Snapshot", page.dom_snapshot_path) +
      "</article>"
    );
  }

  function renderAPICard(api) {
    const title = [api.method || "API", api.path_template || api.url || ""].filter(Boolean).join(" ");
    return (
      '<article class="wb-card">' +
      "<header>" +
      "<div><h3>" +
      escapeHTML(api.semantic_name || title) +
      '</h3><p>' +
      escapeHTML(title) +
      "</p></div>" +
      "</header>" +
      '<div class="wb-card-meta">' +
      pill(escapeHTML(api.operation_type || "unknown")) +
      riskPill(api.risk) +
      (api.trigger_route ? pill(escapeHTML(api.trigger_route)) : "") +
      (api.resource_type ? pill(escapeHTML(api.resource_type)) : "") +
      (api.status ? pill("status " + escapeHTML(String(api.status))) : "") +
      "</div>" +
      (api.trigger_action ? "<p>触发动作：" + escapeHTML(api.trigger_action) + "</p>" : "") +
      (api.content_type ? "<p>Content-Type：" + escapeHTML(api.content_type) + "</p>" : "") +
      (api.url ? renderArtifactPath("Captured URL", api.url, false) : "") +
      renderSchema("Request Schema", api.request_schema) +
      renderSchema("Response Schema", api.response_schema) +
      "</article>"
    );
  }

  function renderEntityCard(entity) {
    const fields = Array.isArray(entity.fields) ? entity.fields : [];
    return (
      '<article class="wb-card">' +
      "<header>" +
      "<div><h3>" +
      escapeHTML(entity.label || entity.name || "未命名实体") +
      '</h3><p>' +
      escapeHTML(entity.name || "") +
      "</p></div>" +
      "</header>" +
      '<div class="wb-card-meta">' +
      fields
        .slice(0, 10)
        .map((field) => pill(escapeHTML(field.label || field.name || "field") + (field.type ? " · " + escapeHTML(field.type) : "")))
        .join("") +
      "</div>" +
      renderSchema("Fields", fields) +
      "</article>"
    );
  }

  function renderOptionalList(label, items) {
    if (!items || !items.length) {
      return "";
    }
    return "<p><strong>" + escapeHTML(label) + "：</strong>" + escapeHTML(items.join("、")) + "</p>";
  }

  function renderSchema(title, value) {
    if (!value) {
      return "";
    }
    return (
      '<details class="wb-schema"><summary>' +
      escapeHTML(title) +
      "</summary><pre>" +
      escapeHTML(formatJSON(value)) +
      "</pre></details>"
    );
  }

  function renderArtifactPath(label, path, useArtifactRoute) {
    const normalized = normalizeString(path);
    if (!normalized) {
      return "";
    }
    const href =
      useArtifactRoute === false ? normalizeArtifactExternalURL(normalized) : artifactHrefForPath(normalized);
    if (href) {
      return (
        '<p><strong>' +
        escapeHTML(label) +
        '：</strong><a href="' +
        escapeAttr(href) +
        '" target="_blank" rel="noreferrer">' +
        escapeHTML(normalized) +
        "</a></p>"
      );
    }
    return "<p><strong>" + escapeHTML(label) + "：</strong>" + escapeHTML(normalized) + "</p>";
  }

  function metaCard(label, value) {
    return (
      '<div class="wb-meta-card"><span>' +
      escapeHTML(label) +
      "</span><strong>" +
      escapeHTML(value) +
      "</strong></div>"
    );
  }

  function pill(text) {
    return '<span class="wb-pill">' + text + "</span>";
  }

  function riskPill(risk) {
    const normalized = normalizeString(risk) || "unknown";
    return '<span class="wb-pill wb-pill-risk-' + escapeHTML(normalized) + '">' + escapeHTML(normalized) + "</span>";
  }

  function buildOptions(items, selectedValue) {
    return items
      .map((item) => {
        const value = normalizeString(item.value);
        const selected = value === normalizeString(selectedValue) ? " selected" : "";
        return '<option value="' + escapeAttr(value) + '"' + selected + ">" + escapeHTML(item.label) + "</option>";
      })
      .join("");
  }

  function findSite(siteId) {
    return state.sites.find((item) => item.site_id === siteId) || null;
  }

  function findProvider(providerId) {
    return state.providers.find((item) => item.provider_id === providerId) || null;
  }

  function describeSite(site) {
    if (!site) {
      return "未选择";
    }
    return [site.name || site.site_id || "", site.site_id || ""].filter(Boolean).join(" · ");
  }

  function describeProviderOption(provider) {
    const pieces = [
      provider.name || provider.provider_id || "",
      provider.type || "",
      provider.resolved_model || provider.model || "",
      provider.ready ? "ready" : provider.status || "",
    ].filter(Boolean);
    return pieces.join(" · ");
  }

  function applyProviderTypeDefaults() {
    const providerType = normalizeString(refs.providerTypeSelect.value);
    if (!normalizeString(refs.providerBaseURLInput.value)) {
      if (providerType === "ollama") {
        refs.providerBaseURLInput.value = "http://127.0.0.1:11434";
      } else {
        refs.providerBaseURLInput.value = "https://api.openai.com/v1";
      }
    }
    if (!normalizeString(refs.providerModelInput.value)) {
      if (providerType === "ollama") {
        refs.providerModelInput.value = "qwen2.5-coder:7b";
      } else {
        refs.providerModelInput.value = "gpt-4.1-mini";
      }
    }
    if (!normalizeString(refs.providerAPIKeyEnvInput.value) && providerType !== "ollama") {
      refs.providerAPIKeyEnvInput.value = "OPENAI_API_KEY";
    }
  }

  function parseList(value) {
    return normalizeString(value)
      .split(/[,\n]/)
      .map((item) => normalizeString(item))
      .filter(Boolean);
  }

  function toNumber(value, fallback) {
    const parsed = Number.parseInt(String(value), 10);
    return Number.isFinite(parsed) ? parsed : fallback;
  }

  function normalizeString(value) {
    return String(value || "").trim();
  }

  function formatJSON(value) {
    try {
      return JSON.stringify(value, null, 2);
    } catch (error) {
      return String(value);
    }
  }

  function formatDate(value) {
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) {
      return value;
    }
    return date.toLocaleString("zh-CN", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    });
  }

  function flattenTraceList(traces) {
    const items = [];
    (traces || []).forEach((trace) => {
      if (!trace) {
        return;
      }
      items.push(trace);
      if (Array.isArray(trace.children) && trace.children.length) {
        items.push.apply(items, flattenTraceList(trace.children));
      }
      if (Array.isArray(trace.attempts) && trace.attempts.length) {
        items.push.apply(items, flattenTraceList(trace.attempts));
      }
      if (trace.condition) {
        items.push.apply(items, flattenTraceList([trace.condition]));
      }
    });
    return items;
  }

  function artifactHrefForPath(path) {
    if (!state.appMeta || !state.appMeta.artifact_base_path) {
      return "";
    }
    const normalizedPath = slashPath(path);
    const artifactRoot = slashPath(state.appMeta.artifact_root || "");
    let relativePath = normalizedPath;
    if (artifactRoot) {
      if (normalizedPath === artifactRoot) {
        relativePath = "";
      } else if (normalizedPath.startsWith(artifactRoot + "/")) {
        relativePath = normalizedPath.slice(artifactRoot.length + 1);
      } else if (normalizedPath.startsWith("/")) {
        return "";
      }
    }
    if (!relativePath) {
      return "";
    }
    return (
      state.appMeta.artifact_base_path.replace(/\/+$/, "/") +
      relativePath
        .split("/")
        .filter(Boolean)
        .map((segment) => encodeURIComponent(segment))
        .join("/")
    );
  }

  function normalizeArtifactExternalURL(path) {
    const value = normalizeString(path);
    if (!value) {
      return "";
    }
    if (/^https?:\/\//i.test(value)) {
      return value;
    }
    return "";
  }

  function slashPath(value) {
    return normalizeString(value).replaceAll("\\", "/").replace(/\/+$/, "");
  }

  function setStatus(message, stateName) {
    refs.statusBar.textContent = message;
    refs.statusBar.dataset.state = stateName || "idle";
  }

  async function apiFetch(path, options) {
    const response = await rawFetch(path, options);
    if (response.status === 404) {
      throw new Error("请求的资源不存在：" + path);
    }
    if (!response.ok) {
      throw new Error(response.error || response.statusText || "请求失败");
    }
    return response.data;
  }

  async function apiFetchOptional(path, fallback) {
    const response = await rawFetch(path);
    if (response.status === 404) {
      return fallback;
    }
    if (!response.ok) {
      throw new Error(response.error || response.statusText || "请求失败");
    }
    return response.data;
  }

  async function rawFetch(path, options) {
    const requestOptions = Object.assign(
      {
        headers: {},
      },
      options || {}
    );
    if (requestOptions.body && !requestOptions.headers["Content-Type"]) {
      requestOptions.headers["Content-Type"] = "application/json";
    }
    const response = await fetch(path, requestOptions);
    const text = await response.text();
    let data = {};
    if (text) {
      try {
        data = JSON.parse(text);
      } catch (error) {
        data = { raw: text };
      }
    }
    return {
      ok: response.ok,
      status: response.status,
      statusText: response.statusText,
      error: data && data.error ? data.error : "",
      data: data,
    };
  }

  function escapeHTML(value) {
    return String(value || "")
      .replaceAll("&", "&amp;")
      .replaceAll("<", "&lt;")
      .replaceAll(">", "&gt;")
      .replaceAll('"', "&quot;")
      .replaceAll("'", "&#39;");
  }

  function escapeAttr(value) {
    return escapeHTML(value).replaceAll("`", "&#96;");
  }
})();
