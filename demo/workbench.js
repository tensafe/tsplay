(function () {
  const state = {
    artifactRoot: "",
    appMeta: null,
    sessions: [],
    providers: [],
    sites: [],
    viewMode: readStoredViewMode(),
    currentSiteId: "",
    currentTab: "pages",
    apiSearchQuery: "",
    pages: [],
    apis: [],
    entities: [],
    lastExplore: null,
    lastPlan: null,
    lastPlanRealtimeContext: null,
    currentFlowYAML: "",
    lastRun: null,
    lastRepair: null,
  };

  const refs = {};

  document.addEventListener("DOMContentLoaded", () => {
    cacheRefs();
    bindEvents();
    renderViewMode();
    renderTabs();
    applyProviderTypeDefaults();
    runAction("加载 Workbench", refreshAll);
  });

  function cacheRefs() {
    refs.currentSiteHeader = document.getElementById("currentSiteHeader");
    refs.currentSiteStatus = document.getElementById("currentSiteStatus");
    refs.currentSiteOverview = document.getElementById("currentSiteOverview");
    refs.modeButtons = Array.from(document.querySelectorAll(".wb-view-button"));
    refs.modeNoviceButton = document.getElementById("modeNoviceButton");
    refs.modeDeveloperButton = document.getElementById("modeDeveloperButton");
    refs.stepSiteCard = document.getElementById("stepSiteCard");
    refs.stepExploreCard = document.getElementById("stepExploreCard");
    refs.stepPlanCard = document.getElementById("stepPlanCard");
    refs.stepSiteCopy = document.getElementById("stepSiteCopy");
    refs.stepExploreCopy = document.getElementById("stepExploreCopy");
    refs.stepPlanCopy = document.getElementById("stepPlanCopy");
    refs.statusBar = document.getElementById("statusBar");
    refs.refreshAllButton = document.getElementById("refreshAllButton");
    refs.openSiteConfigButton = document.getElementById("openSiteConfigButton");
    refs.siteConfigDetails = document.getElementById("siteConfigDetails");
    refs.siteAdvancedDetails = document.getElementById("siteAdvancedDetails");
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
    refs.exploreLifecycle = document.getElementById("exploreLifecycle");
    refs.cardList = document.getElementById("cardList");
    refs.tabs = Array.from(document.querySelectorAll(".wb-tab"));
    refs.apiSearchBar = document.getElementById("apiSearchBar");
    refs.apiSearchInput = document.getElementById("apiSearchInput");
    refs.clearAPISearchButton = document.getElementById("clearAPISearchButton");
    refs.apiSearchMeta = document.getElementById("apiSearchMeta");
    refs.taskForm = document.getElementById("taskForm");
    refs.intentInput = document.getElementById("intentInput");
    refs.planContextPreview = document.getElementById("planContextPreview");
    refs.runHeadlessInput = document.getElementById("runHeadlessInput");
    refs.runFlowButton = document.getElementById("runFlowButton");
    refs.planSummary = document.getElementById("planSummary");
    refs.flowStoryboard = document.getElementById("flowStoryboard");
    refs.planCandidates = document.getElementById("planCandidates");
    refs.originalFlowViewer = document.getElementById("originalFlowViewer");
    refs.flowEditor = document.getElementById("flowEditor");
    refs.copyOriginalFlowButton = document.getElementById("copyOriginalFlowButton");
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
    refs.refreshAllButton.addEventListener("click", () => runAction("刷新数据", refreshAll));
    refs.modeButtons.forEach((button) => {
      button.addEventListener("click", () => {
        setViewMode(button.dataset.viewMode || "novice");
      });
    });
    refs.openSiteConfigButton.addEventListener("click", () => {
      openDetails(refs.siteConfigDetails);
    });
    refs.loadSiteButton.addEventListener("click", () => runAction("加载站点", () => selectSite(refs.sitePicker.value)));
    refs.sitePicker.addEventListener("change", () => runAction("切换站点", () => selectSite(refs.sitePicker.value)));
    refs.siteUrlInput.addEventListener("blur", suggestSiteIdentity);
    refs.siteSessionSelect.addEventListener("change", renderSessionSummary);
    refs.siteProviderSelect.addEventListener("change", renderProviderSummary);
    refs.newSiteButton.addEventListener("click", openNewSiteFlow);
    refs.siteForm.addEventListener("submit", (event) => {
      event.preventDefault();
      runAction("保存站点", saveSite);
    });
    refs.sessionForm.addEventListener("submit", (event) => {
      event.preventDefault();
      runAction("保存 Session", saveSession);
    });
    refs.providerForm.addEventListener("submit", (event) => {
      event.preventDefault();
      runAction("保存 Provider", saveProvider);
    });
    refs.providerTypeSelect.addEventListener("change", applyProviderTypeDefaults);
    refs.exploreForm.addEventListener("submit", (event) => {
      event.preventDefault();
      runAction("探索站点", exploreSite);
    });
    refs.taskForm.addEventListener("submit", (event) => {
      event.preventDefault();
      runAction("生成 Flow", planTask);
    });
    refs.intentInput.addEventListener("input", renderPlanContextPreview);
    refs.flowEditor.addEventListener("input", handleFlowEditorInput);
    refs.runFlowButton.addEventListener("click", () => runAction("执行 Flow", executePlanFlow));
    refs.replayFlowButton.addEventListener("click", () => runAction("回放 Flow", executePlanFlow));
    refs.repairFlowButton.addEventListener("click", () => runAction("生成 Repair Context", buildRepairContext));
    refs.autoRepairFlowButton.addEventListener("click", () => runAction("自动修复 Flow", buildAutoRepair));
    refs.runRepairedFlowButton.addEventListener("click", () => runAction("执行修复版 Flow", executePlanFlow));
    refs.copyOriginalFlowButton.addEventListener("click", () => runAction("复制原始 Flow", copyOriginalFlow));
    refs.copyFlowButton.addEventListener("click", () => runAction("复制 Flow", copyFlow));
    refs.tabs.forEach((button) => {
      button.addEventListener("click", () => {
        state.currentTab = button.dataset.tab || "pages";
        renderTabs();
        renderCards();
      });
    });
    refs.apiSearchInput.addEventListener("input", () => {
      state.apiSearchQuery = normalizeString(refs.apiSearchInput.value);
      renderCards();
    });
    refs.clearAPISearchButton.addEventListener("click", () => {
      state.apiSearchQuery = "";
      refs.apiSearchInput.value = "";
      renderCards();
      refs.apiSearchInput.focus();
    });
    refs.exploreMeta.addEventListener("click", (event) => {
      const target = event.target.closest("[data-tab-target]");
      if (!target) {
        return;
      }
      state.currentTab = target.dataset.tabTarget || "pages";
      renderTabs();
      renderCards();
    });
  }

  async function runAction(label, action) {
    setStatus(describeBusyStatus(label), "busy");
    try {
      const result = await action();
      setStatus(describeSuccessStatus(label, result), "idle");
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

    await Promise.all([loadSessions(), loadProviders(), loadSites()]);
    if (state.currentSiteId) {
      await loadKnowledge(state.currentSiteId);
    } else {
      renderExploreMeta();
      renderExploreLifecycle();
      renderCards();
      renderSiteContext();
    }
    renderJourney();
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
      renderSiteContext();
    }
  }

  async function selectSite(siteId) {
    siteId = normalizeString(siteId);
    state.currentSiteId = siteId;
    state.apiSearchQuery = "";
    refs.apiSearchInput.value = "";
    refs.sitePicker.value = siteId;
    const site = findSite(siteId);
    populateSiteForm(site);
    if (siteId && refs.siteConfigDetails) {
      refs.siteConfigDetails.open = false;
    }
    state.lastPlan = null;
    state.lastPlanRealtimeContext = null;
    state.currentFlowYAML = "";
    state.lastRun = null;
    state.lastRepair = null;
    renderJourney();
    renderPlan();
    renderRunResult();
    renderRepair();
    if (!siteId) {
      state.pages = [];
      state.apis = [];
      state.entities = [];
      state.lastExplore = null;
      renderExploreMeta();
      renderExploreLifecycle();
      renderCards();
      renderJourney();
      renderSiteContext();
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
    const derivedRunID = deriveLatestExploreRunID(state.pages);
    const derivedFinishedAt = deriveLatestKnowledgeTimestamp(state.pages, state.apis, state.entities);
    state.lastExplore = {
      run_id: derivedRunID || (state.lastExplore ? state.lastExplore.run_id : ""),
      explore_mode:
        (state.pages[0] && state.pages[0].discovery_mode) ||
        (state.lastExplore ? state.lastExplore.explore_mode : "") ||
        "",
      pages: state.pages,
      apis: state.apis,
      entities: state.entities,
      finished_at: derivedFinishedAt || (state.lastExplore ? state.lastExplore.finished_at : ""),
    };
    renderSiteContext();
    renderExploreMeta();
    renderExploreLifecycle();
    renderCards();
    renderJourney();
  }

  async function saveSite() {
    suggestSiteIdentity();
    const payload = {
      site_id: normalizeString(refs.siteIdInput.value),
      name: normalizeString(refs.siteNameInput.value),
      start_url: normalizeString(refs.siteUrlInput.value),
      allowed_domains: parseList(refs.siteDomainsInput.value),
      session_name: normalizeString(refs.siteSessionSelect.value),
      provider_id: normalizeString(refs.siteProviderSelect.value),
    };
    if (!payload.start_url) {
      throw new Error("请先填写 Start URL");
    }
    if (!payload.site_id) {
      throw new Error("无法推导站点 ID，请补一个 Start URL 或手动填写站点 ID");
    }
    const saved = await apiFetch("/api/workbench/sites", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    await loadSites(saved.site_id);
    refs.sitePicker.value = saved.site_id;
    await selectSite(saved.site_id);
    return saved;
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
    const siteID = await ensureCurrentSiteReady();
    const baseline = knowledgeSignature();
    const payload = {
      headless: refs.exploreHeadless.checked,
      timeout_ms: toNumber(refs.exploreTimeout.value, 30000),
      max_pages: toNumber(refs.exploreMaxPages.value, 8),
    };
    const path = "/api/workbench/sites/" + encodeURIComponent(siteID) + "/explore";

    let requestDone = false;
    let requestResponse = null;
    let requestError = null;
    rawFetch(path, {
      method: "POST",
      body: JSON.stringify(payload),
    })
      .then((response) => {
        requestDone = true;
        requestResponse = response;
      })
      .catch((error) => {
        requestDone = true;
        requestError = error;
      });

    const deadline = Date.now() + Math.max(payload.timeout_ms + 5000, 15000);
    let statusHintShown = false;
    while (Date.now() < deadline) {
      if (requestDone) {
        if (requestResponse && requestResponse.ok) {
          applyExploreResult(requestResponse.data);
          return;
        }
        if (requestError) {
          throw requestError;
        }
      }

      await sleep(1200);
      await loadKnowledge(siteID);
      if (!statusHintShown) {
        setStatus("探索仍在后台执行，正在自动刷新结果…", "busy");
        statusHintShown = true;
      }
      if (knowledgeSignature() !== baseline) {
        state.lastRepair = null;
        renderRepair();
        return;
      }
    }

    if (requestDone && requestResponse && requestResponse.ok) {
      applyExploreResult(requestResponse.data);
      return;
    }
    throw new Error("探索结果返回较慢，但页面会继续保留已有知识。请稍等后点一次刷新，或重新开始探索。");
  }

  function applyExploreResult(result) {
    state.lastExplore = result;
    state.pages = Array.isArray(result.pages) ? result.pages : [];
    state.apis = Array.isArray(result.apis) ? result.apis : [];
    state.entities = Array.isArray(result.entities) ? result.entities : [];
    state.lastRepair = null;
    renderExploreMeta();
    renderExploreLifecycle();
    renderCards();
    renderJourney();
    renderRepair();
  }

  async function planTask() {
    const siteID = await ensureCurrentSiteReady();
    const intent = normalizeString(refs.intentInput.value);
    if (!intent) {
      throw new Error("请输入任务需求");
    }
    const realtimeContextMeta = await buildRealtimePlanningContext(intent);
    const requestBody = {
      site_id: siteID,
      intent: intent,
      provider_id: normalizeString(refs.siteProviderSelect.value),
    };
    if (realtimeContextMeta && realtimeContextMeta.payload) {
      requestBody.realtime_context = realtimeContextMeta.payload;
    }
    const plan = await apiFetch("/api/workbench/tasks/plan", {
      method: "POST",
      body: JSON.stringify(requestBody),
    });
    state.lastPlan = plan;
    state.lastPlanRealtimeContext = realtimeContextMeta;
    state.currentFlowYAML = typeof plan.flow_yaml === "string" ? plan.flow_yaml : "";
    state.lastRun = null;
    state.lastRepair = null;
    renderJourney();
    renderPlan();
    renderRunResult();
    renderRepair();
  }

  async function executePlanFlow() {
    const siteID = await ensureCurrentSiteReady();
    const flowYAML = normalizeString(currentFlowYAMLValue());
    if (!flowYAML || flowYAML === "还没有生成 Flow") {
      throw new Error("请先生成一个可执行的 Flow");
    }
    const payload = await apiFetch("/api/workbench/tasks/run", {
      method: "POST",
      body: JSON.stringify({
        site_id: siteID,
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
      setCurrentFlowYAML(payload.flow_yaml);
    }
    renderJourney();
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
    const flowYAML = normalizeString(currentFlowYAMLValue());
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
    renderJourney();
    renderRepair();
  }

  async function buildAutoRepair() {
    if (!state.lastRun) {
      throw new Error("当前还没有执行结果，先运行一次 Flow");
    }
    if (state.lastRun.ok) {
      throw new Error("当前执行已成功，不需要自动修复");
    }
    const flowYAML = normalizeString(currentFlowYAMLValue());
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
      setCurrentFlowYAML(payload.repaired_flow_yaml);
    }
    renderJourney();
    renderRepair();
    if (!payload.ok) {
      throw new Error(payload.validation_error || payload.error || "自动修复失败");
    }
  }

  async function copyFlow() {
    const text = normalizeString(currentFlowYAMLValue());
    if (!text) {
      throw new Error("当前没有可复制的 Flow");
    }
    await navigator.clipboard.writeText(text);
  }

  async function copyOriginalFlow() {
    const text = normalizeString(originalFlowYAMLValue());
    if (!text) {
      throw new Error("当前没有原始 Flow 可复制");
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
    renderSessionSummary();
    renderProviderSummary();
    renderSiteContext();
  }

  function resetSiteForm() {
    refs.siteIdInput.value = "";
    refs.siteNameInput.value = "";
    refs.siteUrlInput.value = "";
    refs.siteDomainsInput.value = "";
    refs.siteProviderSelect.value = "";
    refs.sitePicker.value = "";
    state.currentSiteId = "";
    state.apiSearchQuery = "";
    state.pages = [];
    state.apis = [];
    state.entities = [];
    state.lastExplore = null;
    state.lastPlan = null;
    state.currentFlowYAML = "";
    state.lastRun = null;
    state.lastRepair = null;
    renderExploreMeta();
    renderExploreLifecycle();
    renderCards();
    renderJourney();
    renderPlan();
    renderRunResult();
    renderRepair();
    renderProviderSummary();
    renderSiteContext();
  }

  function openNewSiteFlow() {
    resetSiteForm();
    openDetails(refs.siteConfigDetails);
    refs.siteUrlInput.focus();
  }

  function openDetails(element) {
    if (!element) {
      return;
    }
    element.open = true;
  }

  async function ensureCurrentSiteReady() {
    const currentSiteID = normalizeString(state.currentSiteId);
    if (currentSiteID) {
      const knownSite = findSite(currentSiteID);
      if (knownSite) {
        return currentSiteID;
      }
      const startURL = normalizeString(refs.siteUrlInput.value);
      if (!startURL) {
        throw new Error("当前站点配置已经失效，请重新填写 Start URL 后再试");
      }
      setStatus("检测到服务端没有当前站点配置，先自动重新保存站点…", "busy");
      const restored = await saveSite();
      const restoredSiteID = normalizeString(restored && restored.site_id);
      if (!restoredSiteID) {
        throw new Error("自动恢复站点配置失败，请手动点一次保存站点");
      }
      return restoredSiteID;
    }
    const startURL = normalizeString(refs.siteUrlInput.value);
    if (!startURL) {
      throw new Error("请先在第一步填写 Start URL，或者选择一个已保存站点");
    }
    setStatus("检测到当前站点还没保存，先自动保存站点…", "busy");
    const saved = await saveSite();
    const siteID = normalizeString(saved && saved.site_id) || normalizeString(state.currentSiteId);
    if (!siteID) {
      throw new Error("站点保存后仍未拿到 site_id，请重试一次");
    }
    return siteID;
  }

  function suggestSiteIdentity() {
    const rawURL = normalizeString(refs.siteUrlInput.value);
    if (!rawURL) {
      return;
    }
    try {
      const parsed = new URL(rawURL);
      if (!normalizeString(refs.siteNameInput.value)) {
        refs.siteNameInput.value = buildSiteDisplayName(parsed);
      }
      if (!normalizeString(refs.siteIdInput.value)) {
        refs.siteIdInput.value = buildSiteID(parsed);
      }
      if (!normalizeString(refs.siteDomainsInput.value)) {
        refs.siteDomainsInput.value = parsed.hostname;
      }
    } catch (_error) {
      // Ignore invalid URL while typing; required validation still happens on submit.
    }
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
      refs.providerSummary.textContent = "未显式选择 provider 时，系统会优先尝试已就绪的默认 provider 来自动生成或修复 Flow。";
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

  function renderJourney() {
    const hasSite = !!normalizeString(state.currentSiteId);
    const hasExplore = (state.pages && state.pages.length > 0) || (state.apis && state.apis.length > 0);
    const hasPlan = !!(state.lastPlan && normalizeString(state.lastPlan.flow_yaml));

    setStepState(
      refs.stepSiteCard,
      hasSite ? "done" : "active",
      hasSite ? "站点上下文已就绪。" : "先选择一个已保存站点，或者填一个地址保存下来。"
    );

    setStepState(
      refs.stepExploreCard,
      hasExplore ? "done" : hasSite ? "active" : "todo",
      hasExplore
        ? "已拿到页面、接口和实体线索。"
        : hasSite
          ? "执行一次探索，先拿到知识卡片。"
          : "先完成站点接入，再开始探索。"
    );

    setStepState(
      refs.stepPlanCard,
      hasPlan ? "done" : hasExplore ? "active" : "todo",
      hasPlan
        ? "Flow 已起草，先看步骤卡片再执行。"
        : hasExplore
          ? "描述任务目标，系统会先匹配页面和 API。"
          : "站点探索完成后，这一步会更准确。"
    );
  }

  function setStepState(card, status, message) {
    if (!card) {
      return;
    }
    card.classList.toggle("is-active", status === "active");
    card.classList.toggle("is-done", status === "done");
    const text = card.querySelector("p");
    if (text) {
      text.textContent = message;
    }
  }

  function setViewMode(mode) {
    const nextMode = normalizeViewMode(mode);
    if (state.viewMode === nextMode) {
      renderViewMode();
      return;
    }
    state.viewMode = nextMode;
    writeStoredViewMode(nextMode);
    if (nextMode !== "developer" && refs.siteAdvancedDetails) {
      refs.siteAdvancedDetails.open = false;
    }
    renderViewMode();
    renderExploreLifecycle();
    renderCards();
    renderPlan();
    renderRunResult();
    renderRepair();
    setStatus(nextMode === "developer" ? "已展开开发者详情" : "已切换到工作台视图", "idle");
  }

  function renderViewMode() {
    document.body.dataset.viewMode = state.viewMode;
    refs.modeButtons.forEach((button) => {
      button.classList.toggle("is-active", normalizeViewMode(button.dataset.viewMode) === state.viewMode);
    });
    renderKnowledgeToolbar(state.currentTab === "apis" ? state.apis.length : 0, state.currentTab === "apis" ? filterAPIs(state.apis).length : 0);
  }

  function renderTabs() {
    const counts = {
      pages: state.pages.length,
      apis: state.apis.length,
      entities: state.entities.length,
    };
    refs.tabs.forEach((button) => {
      const tabName = button.dataset.tab || "pages";
      const label = button.dataset.label || button.textContent || "";
      button.classList.toggle("is-active", tabName === state.currentTab);
      button.innerHTML =
        escapeHTML(label) +
        ' <span class="wb-tab-count">' +
        escapeHTML(String(counts[tabName] || 0)) +
        "</span>";
    });
    renderKnowledgeToolbar(counts.apis, filterAPIs(state.apis).length);
  }

  function renderExploreMeta() {
    const eventCount = totalPageEventCount(state.pages);
    refs.exploreMeta.innerHTML = [
      metaCard("页面", String(state.pages.length), "pages"),
      metaCard("API", String(state.apis.length), "apis"),
      metaCard("实体", String(state.entities.length), "entities"),
      metaCard("事件", String(eventCount), "pages"),
    ].join("");
  }

  function renderExploreLifecycle() {
    if (!refs.exploreLifecycle) {
      return;
    }
    const hasKnowledge = state.pages.length || state.apis.length || state.entities.length;
    if (!hasKnowledge) {
      refs.exploreLifecycle.innerHTML =
        "<h3>运行结论</h3><p>开始探索后，这里会先展示页面、API、实体和事件摘要，详细请求和运行信息放进“运行详情”。</p>";
      return;
    }

    const summary = summarizeExploreLifecycle(state.pages);
    const runID = state.lastExplore && state.lastExplore.run_id ? state.lastExplore.run_id : "暂无";
    const finishedAt = state.lastExplore && state.lastExplore.finished_at ? formatDate(state.lastExplore.finished_at) : "暂无";
    const exploreMode = state.lastExplore && state.lastExplore.explore_mode
      ? describeExploreMode(state.lastExplore.explore_mode)
      : "暂无";
    const observationStatus = describeObservationStatus(summary);
    refs.exploreLifecycle.innerHTML =
      "<h3>运行结论</h3>" +
      "<p>" +
      escapeHTML(describeExploreHeadline()) +
      "</p>" +
      '<div class="wb-card-meta">' +
      pill("模式 · " + escapeHTML(exploreMode)) +
      "</div>" +
      (observationStatus.message ? '<p class="wb-data-note">' + escapeHTML(observationStatus.message) + "</p>" : "") +
      '<details class="wb-details wb-details-subtle">' +
      "<summary>运行详情</summary>" +
      '<div class="wb-detail-grid">' +
      miniStat("最近探索", finishedAt) +
      miniStat("请求数", String(summary.networkRequestCount)) +
      miniStat("可读响应数", String(summary.readableResponseCount)) +
      miniStat("失败请求数", String(summary.networkFailureCount)) +
      miniStat("观察错误数", String(summary.observationErrorCount)) +
      miniStat("交互元素", String(summary.interactiveElementCount)) +
      miniStat("内容块", String(summary.contentElementCount)) +
      (isDeveloperMode()
        ? miniStat("Run ID", runID) +
          miniStat("Artifact Root", state.appMeta && state.appMeta.artifact_root ? state.appMeta.artifact_root : state.artifactRoot || "default")
        : "") +
      "</div>" +
      '<p class="wb-data-note"><strong>加载期回调：</strong>' +
      escapeHTML(summary.filterRule) +
      "</p>" +
      '<p class="wb-data-note"><strong>渲染后观察：</strong>' +
      escapeHTML(summary.observationMode) +
      "</p>" +
      (isDeveloperMode() && observationStatus.technical
        ? renderDisclosure("查看技术详情", '<p class="wb-data-note">' + escapeHTML(observationStatus.technical) + "</p>")
        : "") +
      "</details>";
  }

  function renderCards() {
    renderTabs();
    const items =
      state.currentTab === "pages"
        ? state.pages
        : state.currentTab === "apis"
          ? filterAPIs(state.apis)
          : state.entities;
    if (!items.length) {
      refs.cardList.innerHTML = renderEmptyCardForCurrentTab();
      renderPlanContextPreview();
      return;
    }
    if (state.currentTab === "apis") {
      refs.cardList.innerHTML = renderAPIList(items);
      renderPlanContextPreview();
      return;
    }
    const renderers = {
      pages: renderPageCard,
      entities: renderEntityCard,
    };
    refs.cardList.innerHTML = items.map((item) => renderers[state.currentTab](item)).join("");
    renderPlanContextPreview();
  }

  function renderPlan() {
    const plan = state.lastPlan;
    if (!plan) {
      refs.planSummary.innerHTML =
        "<h3>规划结果</h3><p>先输入任务需求，系统会根据当前站点自动匹配页面和 API，再生成步骤卡片。</p>";
      refs.flowStoryboard.innerHTML = "";
      refs.planCandidates.innerHTML = "";
      refs.originalFlowViewer.value = "# 还没有生成原始 Flow";
      setCurrentFlowYAML("");
      refs.autoRepairFlowButton.disabled = true;
      return;
    }

    const provider = plan.provider || null;
    const realtimeContext = state.lastPlanRealtimeContext;
    const warnings = Array.isArray(plan.warnings) ? plan.warnings : [];
    const summaryBits = [
      '<div class="wb-card-meta">',
      pill("策略 · " + escapeHTML(plan.strategy || "unknown")),
      pill("生成 · " + escapeHTML(plan.generation_mode || "local")),
      pill("Flow · " + escapeHTML(plan.flow_name || "未命名")),
      (provider ? pill("Provider · " + escapeHTML(provider.name || provider.provider_id || "auto")) : ""),
      "</div>",
      "<p>" + escapeHTML(plan.reason || "系统已经生成了一条可执行的建议路径。") + "</p>",
      realtimeContext ? "<p><strong>规划上下文：</strong>" + escapeHTML(describeRealtimeContextMeta(realtimeContext)) + "</p>" : "",
      realtimeContext && realtimeContext.warnings.length
        ? "<p><strong>上下文提示：</strong>" + escapeHTML(realtimeContext.warnings.join("；")) + "</p>"
        : "",
      plan.validation_error ? "<p><strong>AI 校验失败：</strong>" + escapeHTML(plan.validation_error) + "</p>" : "",
      warnings.length ? "<p><strong>提示：</strong>" + escapeHTML(warnings.join("；")) + "</p>" : "",
    ];
    refs.planSummary.innerHTML = "<h3>规划结果</h3>" + summaryBits.join("");
    refs.flowStoryboard.innerHTML = renderFlowStoryboard(plan.flow);

    const sections = [];
    if (Array.isArray(plan.matched_pages) && plan.matched_pages.length) {
      sections.push(renderCandidateSection("命中的页面", plan.matched_pages));
    }
    if (Array.isArray(plan.matched_apis) && plan.matched_apis.length) {
      sections.push(renderCandidateSection("命中的 API", plan.matched_apis));
    }
    if (isDeveloperMode() && realtimeContext && realtimeContext.payload) {
      sections.push(renderSchema("Realtime Context", realtimeContext.payload));
    }
    if (isDeveloperMode() && plan.model_output) {
      sections.push(renderSchema("Model Output", plan.model_output));
    }
    refs.planCandidates.innerHTML = sections.join("");
    refs.originalFlowViewer.value = originalFlowYAMLValue() || "# 当前没有原始 Flow";
    if (!currentFlowYAMLValue()) {
      setCurrentFlowYAML(plan.flow_yaml || "");
    } else {
      refs.flowEditor.value = currentFlowYAMLValue();
      syncFlowActionButtons();
    }
    refs.autoRepairFlowButton.disabled = true;
  }

  function renderRunResult() {
    const run = state.lastRun;
    if (!run) {
      refs.runResultPanel.innerHTML =
        "<h3>执行结果</h3><p>生成 Flow 后，可以直接在这里执行，并查看当前执行结果摘要。</p>";
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
      (isDeveloperMode() && result.trace ? pill("Trace · " + escapeHTML(String(result.trace.length || 0))) : ""),
      "</div>",
      run.error ? "<p>" + escapeHTML(run.error) + "</p>" : "<p>Flow 已执行完成，可以继续下一步。</p>",
      isDeveloperMode() ? renderArtifactPath("Run Root", result.run_root) : "",
      isDeveloperMode() ? renderArtifactPath("Browser Video", result.browser_video) : "",
      isDeveloperMode() ? renderSchema("Vars", result.vars) : "",
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

  function renderFlowStoryboard(flow) {
    const steps = flow && Array.isArray(flow.steps) ? flow.steps : [];
    if (!steps.length) {
      return '<article class="wb-empty-card"><h3>还没有步骤卡片</h3><p>当前 Flow 还没有可展示的步骤。</p></article>';
    }
    return (
      '<div class="wb-story-list">' +
      steps
        .map((step, index) => {
          const summary = describeFlowStep(step, index);
          return (
            '<article class="wb-story-step">' +
            '<span class="wb-story-badge">' + escapeHTML(String(index + 1)) + "</span>" +
            "<div><strong>" +
            escapeHTML(summary.title) +
            "</strong><p>" +
            escapeHTML(summary.body) +
            "</p></div></article>"
          );
        })
        .join("") +
      "</div>"
    );
  }

  function renderPlanContextPreview() {
    if (!refs.planContextPreview) {
      return;
    }
    const intent = normalizeString(refs.intentInput ? refs.intentInput.value : "");
    const page = pickRealtimeContextPage(intent);
    if (!state.currentSiteId) {
      refs.planContextPreview.textContent = "先选择一个站点，生成 Flow 时系统才会自动附带该站点的实时上下文。";
      return;
    }
    if (!page) {
      refs.planContextPreview.textContent =
        state.pages.length > 0
          ? "当前页面卡片里还没有足够稳定的上下文，生成 Flow 时会先退回到站点知识库。"
          : "当前还没有页面卡片；先探索一次站点，生成 Flow 时才能自动附带 Observation。";
      return;
    }

    const pageLabel = firstNonEmpty(
      normalizeString(page.title),
      normalizeString(page.normalized_route),
      normalizeString(page.url),
      "页面卡片"
    );
    const details = [];
    const url = normalizeString(page.url);
    if (url) {
      details.push(url);
    }
    details.push(page.observation_path ? "会附带 Observation" : "只会附带 URL / 标题");
    refs.planContextPreview.innerHTML =
      "<strong>本次规划会优先使用：</strong>" +
      escapeHTML(pageLabel) +
      (details.length ? '<p class="wb-data-note">' + escapeHTML(details.join(" · ")) + "</p>" : "");
  }

  async function buildRealtimePlanningContext(intent) {
    const page = pickRealtimeContextPage(intent);
    if (!page) {
      return null;
    }

    const payload = {};
    const warnings = [];
    const pageURL = normalizeString(page.url);
    const pageTitle = normalizeString(page.title);
    const observationPath = normalizeString(page.observation_path);

    if (pageURL) {
      payload.url = pageURL;
    }
    if (pageTitle) {
      payload.title = pageTitle;
    }

    if (observationPath) {
      try {
        const observation = await loadArtifactJSON(observationPath);
        if (hasMeaningfulRealtimeObservation(observation)) {
          payload.observation = observation;
          payload.url = firstNonEmpty(normalizeString(payload.url), normalizeString(observation.url));
          payload.title = firstNonEmpty(normalizeString(payload.title), normalizeString(observation.title));
        } else {
          warnings.push("Observation artifact 为空，已退回到页面 URL 和标题。");
        }
      } catch (error) {
        warnings.push("Observation artifact 读取失败，已退回到页面 URL 和标题。");
      }
    }

    if (!payload.url && !payload.title && !payload.observation) {
      return null;
    }

    return {
      page_label: firstNonEmpty(pageTitle, normalizeString(page.normalized_route), pageURL, "页面卡片"),
      observation_path: observationPath,
      payload: payload,
      warnings: warnings,
    };
  }

  function pickRealtimeContextPage(intent) {
    const pages = Array.isArray(state.pages) ? state.pages.filter(Boolean) : [];
    if (!pages.length) {
      return null;
    }
    const site = findSite(state.currentSiteId);
    let bestPage = null;
    let bestScore = Number.NEGATIVE_INFINITY;
    pages.forEach((page, index) => {
      const score = scoreRealtimeContextPage(page, intent, site, index);
      if (!bestPage || score > bestScore) {
        bestPage = page;
        bestScore = score;
      }
    });
    return bestPage;
  }

  function scoreRealtimeContextPage(page, intent, site, index) {
    let score = 0;
    const title = normalizeString(page && page.title);
    const url = normalizeString(page && page.url);
    const route = normalizeString(page && page.normalized_route);
    const summary = normalizeString(page && page.summary);
    const snippets = Array.isArray(page && page.text_snippets) ? page.text_snippets.filter(Boolean) : [];
    const haystack = [title, url, route, summary]
      .concat(snippets)
      .join(" ")
      .toLowerCase();

    if (normalizeString(page && page.observation_path)) {
      score += 50;
    }
    if (page && page.capture_summary) {
      score += Math.min(20, toNumber(page.capture_summary.interactive_element_count, 0));
      score += Math.min(20, toNumber(page.capture_summary.content_element_count, 0));
    }

    const siteStartURL = normalizeString(site && site.start_url);
    if (siteStartURL && samePlanningURL(url, siteStartURL)) {
      score += 120;
    }
    if (route === "/" || route.endsWith(":/")) {
      score += 40;
    }
    if (index === 0) {
      score += 10;
    }

    buildPlanningIntentTerms(intent).forEach((term) => {
      if (!term || !haystack.includes(term)) {
        return;
      }
      score += term.length >= 4 ? 18 : 10;
      if (title.toLowerCase().includes(term)) {
        score += 8;
      }
      if (summary.toLowerCase().includes(term)) {
        score += 6;
      }
    });
    return score;
  }

  function buildPlanningIntentTerms(intent) {
    const value = normalizeString(intent).toLowerCase();
    if (!value) {
      return [];
    }
    const terms = new Set();
    const asciiTerms = value.match(/[a-z0-9_/-]{2,}/g) || [];
    asciiTerms.forEach((term) => terms.add(term));

    const cjkGroups = value.match(/[\u4e00-\u9fff]{2,}/g) || [];
    cjkGroups.forEach((group) => {
      if (group.length <= 4) {
        terms.add(group);
      }
      const maxWindow = Math.min(4, group.length);
      for (let size = 2; size <= maxWindow; size += 1) {
        for (let start = 0; start <= group.length - size; start += 1) {
          terms.add(group.slice(start, start + size));
        }
      }
    });

    return Array.from(terms)
      .map((term) => term.trim())
      .filter((term) => term.length >= 2);
  }

  function samePlanningURL(left, right) {
    const normalizedLeft = normalizePlanningURL(left);
    const normalizedRight = normalizePlanningURL(right);
    return !!normalizedLeft && normalizedLeft === normalizedRight;
  }

  function normalizePlanningURL(value) {
    const input = normalizeString(value);
    if (!input) {
      return "";
    }
    try {
      const parsed = new URL(input);
      return parsed.origin.toLowerCase() + parsed.pathname.replace(/\/+$/, "");
    } catch (error) {
      return input.replace(/\/+$/, "").toLowerCase();
    }
  }

  function describeRealtimeContextMeta(meta) {
    if (!meta || !meta.payload) {
      return "未附带实时上下文";
    }
    const pieces = [];
    if (meta.page_label) {
      pieces.push(meta.page_label);
    }
    if (meta.payload.url) {
      pieces.push(meta.payload.url);
    }
    pieces.push(meta.payload.observation ? "已附带 Observation" : "未附带 Observation");
    return pieces.join(" · ");
  }

  function hasMeaningfulRealtimeObservation(observation) {
    if (!observation || typeof observation !== "object") {
      return false;
    }
    return !!(
      normalizeString(observation.url) ||
      normalizeString(observation.title) ||
      normalizeString(observation.page_summary) ||
      normalizeString(observation.dom_snapshot_excerpt) ||
      (Array.isArray(observation.elements) && observation.elements.length) ||
      (Array.isArray(observation.content_elements) && observation.content_elements.length)
    );
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

  function describeFlowStep(step, index) {
    const action = normalizeString(step && step.action);
    const selector = normalizeString(step && step.selector);
    const url = normalizeString(step && step.url);
    const text = normalizeString(step && step.text);
    const saveAs = normalizeString(step && step.save_as);
    const method = normalizeString(step && step.method);

    switch (action) {
      case "navigate":
        return {
          title: "打开页面",
          body: url ? "进入 " + url : "进入目标页面并等待页面可用。",
        };
      case "wait_for":
      case "wait_for_selector":
        return {
          title: "等待页面状态",
          body: selector ? "等待页面上出现 " + selector + " 后再继续。" : "等待页面加载到可继续执行的状态。",
        };
      case "click":
        return {
          title: "点击页面元素",
          body: selector ? "尝试点击 " + selector + "，推进到下一步。" : "点击当前目标元素，推进任务。",
        };
      case "fill":
      case "type":
        return {
          title: "填写输入内容",
          body: selector ? "向 " + selector + " 写入任务所需内容。" : "填写本步需要的输入内容。",
        };
      case "http_request":
        return {
          title: "直接请求接口",
          body: [method || "GET", url || "目标接口"].filter(Boolean).join(" ") + "，优先直接拿结构化数据。",
        };
      case "set_var":
        return {
          title: "准备执行变量",
          body: saveAs ? "为后续步骤保存变量 " + saveAs + "。" : "准备后续步骤要用到的变量。",
        };
      case "assert_text":
        return {
          title: "确认页面结果",
          body: text ? "检查页面中是否出现 “" + text + "”。" : "检查页面结果是否符合预期。",
        };
      default:
        return {
          title: step && step.name ? step.name : "执行步骤 " + (index + 1),
          body: action ? "动作类型：" + action + (selector ? "，目标：" + selector : "") : "执行一条系统生成的 Flow 步骤。",
        };
    }
  }

  function buildSiteDisplayName(parsedURL) {
    const host = normalizeString(parsedURL.hostname).replace(/^www\./, "");
    const path = normalizeString(parsedURL.pathname)
      .replace(/\/+/g, " ")
      .replace(/[^\w\u4e00-\u9fa5\s-]/g, " ")
      .trim();
    return [host, path].filter(Boolean).join(" · ");
  }

  function buildSiteID(parsedURL) {
    const host = normalizeString(parsedURL.hostname)
      .replace(/^www\./, "")
      .replace(/[^\w]+/g, "_")
      .replace(/^_+|_+$/g, "");
    const path = normalizeString(parsedURL.pathname)
      .replace(/[^\w]+/g, "_")
      .replace(/^_+|_+$/g, "");
    return [host, path].filter(Boolean).join("_") || "site";
  }

  function renderPageCard(page) {
    const routeLabel = page.normalized_route || page.url || "未命名页面";
    const modeLabel = describeExploreMode(page.discovery_mode);
    const actionLabels = Array.isArray(page.actions)
      ? page.actions.map((item) => item.label || item.name || item.selector).filter(Boolean)
      : [];
    const tableLabels = Array.isArray(page.tables)
      ? page.tables.map((item) => item.name || item.selector).filter(Boolean)
      : [];
    const formLabels = Array.isArray(page.forms)
      ? page.forms.map((item) => {
          const fields = Array.isArray(item.fields) ? item.fields.map((field) => field.label || field.name).filter(Boolean) : [];
          return [item.name || item.selector || "表单", fields.join(" / ")].filter(Boolean).join(" · ");
        })
      : [];
    const inputLabels = Array.isArray(page.input_fields)
      ? page.input_fields
          .map((field) => [field.label || field.name || "input", field.selector || ""].filter(Boolean).join(" · "))
          .filter(Boolean)
      : [];
    const keyElements = Array.isArray(page.key_elements)
      ? page.key_elements
          .map((item) => {
            if (typeof item === "string") {
              return item;
            }
            return item && (item.label || item.text || item.name || item.selector || item.role || "");
          })
          .filter(Boolean)
          .slice(0, 8)
      : [];
    const linkGroups = Array.isArray(page.link_groups)
      ? page.link_groups
          .map((item) => {
            if (typeof item === "string") {
              return item;
            }
            return item && (item.label || item.group || item.name || item.selector || "");
          })
          .filter(Boolean)
          .slice(0, 4)
      : [];
    const textLabels = Array.isArray(page.text_snippets) ? page.text_snippets.slice(0, 6) : [];
    const flowHints = page.flow_hints || null;
    const links = Array.isArray(page.links) ? page.links : [];
    const apiHits = Array.isArray(page.api_hits) ? page.api_hits : [];
    const events = Array.isArray(page.events) ? page.events : [];
    const captureSummary = page.capture_summary || null;
    const observationStatus = describeObservationStatus(captureSummary);
    return (
      '<article class="wb-card">' +
      "<header>" +
      "<div>" +
      '<span class="wb-card-kicker">' + escapeHTML(routeLabel) + "</span>" +
      "<h3>" + escapeHTML(page.title || routeLabel) + "</h3>" +
      '<p class="wb-card-subtitle">' + escapeHTML(modeLabel + " · " + formatRiskLabel(page.risk)) + "</p>" +
      "</div>" +
      "</header>" +
      '<p class="wb-card-summary">识别到 ' +
      escapeHTML(String(apiHits.length)) +
      " 个接口、" +
      escapeHTML(String(events.length)) +
      " 条页面事件。</p>" +
      '<p class="wb-data-note">' + escapeHTML(observationStatus.message) + "</p>" +
      renderDisclosure(
        "查看 API",
        renderAPIPreview(apiHits),
        "当前页面还没有识别到接口请求。"
      ) +
      renderDisclosure(
        "查看事件",
        renderEventPreview(events),
        "当前页面还没有页面事件。"
      ) +
      renderDisclosure(
        "查看 Observation",
        renderObservationDetails(page, {
          formLabels: formLabels,
          inputLabels: inputLabels,
          keyElements: keyElements,
          actionLabels: actionLabels,
          tableLabels: tableLabels,
          linkGroups: linkGroups,
          flowHints: flowHints,
          textLabels: textLabels,
        }),
        "当前页面还没有 Observation 证据。"
      ) +
      "</article>"
    );
  }

  function renderCaptureSummary(summary) {
    if (!summary) {
      return "";
    }
    const observationStatus = describeObservationStatus(summary);
    return (
      '<div class="wb-detail-grid">' +
      miniStat("请求数", String(summary.network_request_count || 0)) +
      miniStat("结构化响应", String(summary.readable_response_count || 0)) +
      miniStat("失败请求", String(summary.network_failure_count || 0)) +
      miniStat("页面事件", String(summary.event_count || 0)) +
      miniStat("交互元素", String(summary.interactive_element_count || 0)) +
      miniStat("内容块", String(summary.content_element_count || 0)) +
      miniStat("观察错误", String(summary.observation_error_count || 0)) +
      "</div>" +
      '<p class="wb-data-note">' + escapeHTML(observationStatus.message) + "</p>" +
      '<p class="wb-data-note"><strong>加载期回调：</strong>' +
      escapeHTML(summary.filter_rule || "仅保留 xhr/fetch 网络请求。") +
      "</p>" +
      '<p class="wb-data-note"><strong>渲染后观察：</strong>' +
      escapeHTML(summary.observation_mode || "页面加载完成后观察渲染态 DOM。") +
      "</p>" +
      (isDeveloperMode() && observationStatus.technical
        ? renderDisclosure("查看技术详情", '<p class="wb-data-note">' + escapeHTML(observationStatus.technical) + "</p>")
        : "")
    );
  }

  function renderAPIList(apis) {
    return (
      '<div class="wb-api-table">' +
      '<div class="wb-api-header">' +
      '<span>Method</span>' +
      '<span>Path</span>' +
      '<span>Status</span>' +
      '<span>Type</span>' +
      '<span>Risk</span>' +
      "</div>" +
      apis.map((api) => renderAPIRow(api)).join("") +
      "</div>"
    );
  }

  function renderAPIRow(api) {
    const status = api.status ? String(api.status) : "—";
    const typeLabel = api.operation_type || describePayloadType(api.content_type, api.resource_type) || "—";
    const pathLabel = api.path_template || api.url || "暂无路径";
    return (
      '<details class="wb-api-row">' +
      "<summary>" +
      '<span class="wb-api-cell wb-api-cell-method" data-label="Method">' + escapeHTML(api.method || "GET") + "</span>" +
      '<span class="wb-api-cell wb-api-cell-path" data-label="Path">' + escapeHTML(pathLabel) + "</span>" +
      '<span class="wb-api-cell" data-label="Status">' + escapeHTML(status) + "</span>" +
      '<span class="wb-api-cell" data-label="Type">' + escapeHTML(typeLabel) + "</span>" +
      '<span class="wb-api-cell wb-api-cell-risk" data-label="Risk">' + riskPill(api.risk) + "</span>" +
      "</summary>" +
      '<div class="wb-api-detail">' +
      '<p class="wb-data-note"><strong>请求方法：</strong>' + escapeHTML(api.method || "GET") + "</p>" +
      '<p class="wb-data-note"><strong>路径模板：</strong>' + escapeHTML(pathLabel) + "</p>" +
      (api.url ? '<p class="wb-data-note"><strong>请求 URL：</strong>' + escapeHTML(api.url) + "</p>" : "") +
      (api.content_type ? '<p class="wb-data-note"><strong>Content-Type：</strong>' + escapeHTML(api.content_type) + "</p>" : "") +
      (api.trigger_route ? '<p class="wb-data-note"><strong>来源页面：</strong>' + escapeHTML(api.trigger_route) + "</p>" : "") +
      (api.trigger_action ? '<p class="wb-data-note"><strong>触发动作：</strong>' + escapeHTML(api.trigger_action) + "</p>" : "") +
      (api.error ? '<p class="wb-data-note"><strong>错误信息：</strong>' + escapeHTML(api.error) + "</p>" : "") +
      renderSchema("Request Schema", api.request_schema) +
      renderSchema("Response Schema", api.response_schema) +
      "</div>" +
      "</details>"
    );
  }

  function renderEntityCard(entity) {
    const fields = Array.isArray(entity.fields) ? entity.fields : [];
    return (
      '<details class="wb-card wb-card-collapsible">' +
      "<summary>" +
      '<div class="wb-data-title"><div><strong>' +
      escapeHTML(entity.label || entity.name || "未命名实体") +
      "</strong>" +
      '<p class="wb-data-meta">' +
      escapeHTML(entity.name || "字段摘要") +
      "</p></div>" +
      '<div class="wb-card-meta">' +
      pill("字段 " + escapeHTML(String(fields.length))) +
      "</div></div>" +
      "</summary>" +
      '<div class="wb-card-body"><div class="wb-card-meta">' +
      fields
        .slice(0, 10)
        .map((field) => pill(escapeHTML(field.label || field.name || "field") + (field.type ? " · " + escapeHTML(field.type) : "")))
        .join("") +
      "</div>" +
      renderSchema("Fields", fields) +
      "</div></details>"
    );
  }

  function renderOptionalList(label, items) {
    if (!items || !items.length) {
      return "";
    }
    return "<p><strong>" + escapeHTML(label) + "：</strong>" + escapeHTML(items.join("、")) + "</p>";
  }

  function renderFlowHints(flowHints) {
    if (!flowHints) {
      return "";
    }
    const groups = [
      { title: "推荐主输入", items: normalizeHintItems(flowHints.primary_inputs) },
      { title: "推荐主动作", items: normalizeHintItems(flowHints.primary_actions) },
      { title: "推荐等待条件", items: normalizeHintItems(flowHints.wait_conditions) },
      { title: "推荐 API 优先级", items: normalizeHintItems(flowHints.api_priority) },
    ].filter((group) => group.items.length);
    if (!groups.length) {
      return "";
    }
    return (
      '<div class="wb-flow-hints">' +
      groups
        .map((group) => {
          return (
            '<section class="wb-flow-hint">' +
            "<strong>" +
            escapeHTML(group.title) +
            "</strong><p>" +
            escapeHTML(group.items.join("、")) +
            "</p></section>"
          );
        })
        .join("") +
      "</div>"
    );
  }

  function normalizeHintItems(items) {
    return (Array.isArray(items) ? items : [])
      .map((item) => {
        if (typeof item === "string") {
          return item;
        }
        return item && (item.label || item.text || item.name || item.selector || item.path || "");
      })
      .filter(Boolean);
  }

  function renderLinkPreview(links) {
    const items = Array.isArray(links) ? links.slice(0, 8) : [];
    if (!items.length) {
      return "";
    }
    const hiddenCount = Math.max(0, (Array.isArray(links) ? links.length : 0) - items.length);
    return (
      '<div class="wb-link-preview">' +
      "<strong>链接：</strong>" +
      '<div class="wb-link-list">' +
      items
        .map((item) => {
          const href = normalizeArtifactExternalURL(item.href || "");
          const label = item.text || item.href || "link";
          return (
            '<a class="wb-link-item" href="' +
            escapeAttr(href || "#") +
            '" target="_blank" rel="noreferrer">' +
            escapeHTML(label) +
            "</a>"
          );
        })
        .join("") +
      "</div>" +
      (hiddenCount > 0 ? '<p class="wb-link-more">其余 ' + escapeHTML(String(hiddenCount)) + " 个链接已折叠。</p>" : "") +
      "</div>"
    );
  }

  function renderEventPreview(events) {
    const items = Array.isArray(events) ? events : [];
    if (!items.length) {
      return "";
    }
    const levels = summarizeEventLevels(items);
    return (
      '<p class="wb-data-note">页面事件：' +
      escapeHTML(String(items.length)) +
      " 条</p>" +
      '<p class="wb-data-note">warning ' +
      escapeHTML(String(levels.warning)) +
      " · info " +
      escapeHTML(String(levels.info)) +
      " · error " +
      escapeHTML(String(levels.error)) +
      "。</p>" +
      '<div class="wb-event-list">' +
      items
        .map((item) => {
          const meta = [describePageEventType(item.type), item.level || ""].filter(Boolean).join(" · ");
          const detail = [item.message || "", item.detail || "", item.url || ""].filter(Boolean).join(" · ");
          return (
            '<article class="wb-event-item">' +
            '<div class="wb-card-meta">' +
            pill(escapeHTML(meta || "event")) +
            "</div>" +
            "<p>" +
            escapeHTML(detail || "捕获到页面事件") +
            "</p>" +
            "</article>"
          );
        })
        .join("") +
      "</div>"
    );
  }

  function renderAPIPreview(apiHits) {
    const items = Array.isArray(apiHits) ? apiHits : [];
    if (!items.length) {
      return "";
    }
    return (
      '<div class="wb-data-list">' +
      items
        .map((item) => {
          const title = [item.method || "API", item.path_template || item.url || ""].filter(Boolean).join(" ");
          const meta = [
            item.status ? String(item.status) : "",
            item.operation_type || "",
            describePayloadType(item.content_type, item.resource_type),
          ]
            .filter(Boolean)
            .join(" · ");
          return (
            '<article class="wb-data-row">' +
            '<div class="wb-data-title"><div><strong>' +
            escapeHTML(title) +
            "</strong>" +
            '<p class="wb-data-meta">' +
            escapeHTML(meta || "接口请求") +
            "</p></div>" +
            '<div class="wb-card-meta">' +
            riskPill(item.risk) +
            "</div></div>" +
            (item.content_type ? '<p class="wb-data-note">' + escapeHTML(item.content_type) + "</p>" : "") +
            (item.error ? '<p class="wb-data-note">' + escapeHTML(item.error) + "</p>" : "") +
            "</article>"
          );
        })
        .join("") +
      "</div>"
    );
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

  function renderObservationDetails(page, details) {
    const observationStatus = describeObservationStatus(page.capture_summary);
    const blocks = [
      observationStatus.message ? '<p class="wb-data-note">' + escapeHTML(observationStatus.message) + "</p>" : "",
      renderCaptureSummary(page.capture_summary),
      renderOptionalList("面包屑", page.breadcrumbs || []),
      renderOptionalList("表单", details.formLabels || []),
      renderOptionalList("输入控件", details.inputLabels || []),
      renderOptionalList("关键元素", details.keyElements || []),
      renderOptionalList("动作", details.actionLabels || []),
      renderOptionalList("表格", details.tableLabels || []),
      renderOptionalList("链接分组", details.linkGroups || []),
      renderOptionalList("页面文本", details.textLabels || []),
      renderFlowHints(details.flowHints || null),
      renderLinkPreview(page.links || []),
      isDeveloperMode() ? renderArtifactPath("Observation", page.observation_path) : "",
      isDeveloperMode() ? renderArtifactPath("Screenshot", page.screenshot_path) : "",
      isDeveloperMode() ? renderArtifactPath("DOM Snapshot", page.dom_snapshot_path) : "",
      isDeveloperMode() && observationStatus.technical
        ? renderDisclosure("查看技术详情", '<p class="wb-data-note">' + escapeHTML(observationStatus.technical) + "</p>")
        : "",
    ].filter(Boolean);
    return blocks.join("");
  }

  function renderDisclosure(title, content, emptyText) {
    return (
      '<details class="wb-details wb-details-subtle">' +
      "<summary>" +
      escapeHTML(title) +
      "</summary>" +
      (content || '<p class="wb-data-note">' + escapeHTML(emptyText || "暂无内容。") + "</p>") +
      "</details>"
    );
  }

  function renderKnowledgeToolbar(totalCount, filteredCount) {
    if (!refs.apiSearchBar) {
      return;
    }
    const isAPITab = state.currentTab === "apis";
    refs.apiSearchBar.classList.toggle("is-hidden", !isAPITab);
    if (!isAPITab) {
      refs.apiSearchMeta.textContent = "";
      return;
    }
    refs.apiSearchInput.value = state.apiSearchQuery;
    refs.clearAPISearchButton.disabled = !state.apiSearchQuery;
    if (!totalCount) {
      refs.apiSearchMeta.textContent = "当前还没有 API 卡片。";
      return;
    }
    refs.apiSearchMeta.textContent =
      filteredCount === totalCount
        ? "共 " + totalCount + " 个 API。"
        : "共 " + totalCount + " 个 API，当前命中 " + filteredCount + " 个。";
  }

  function filterAPIs(apis) {
    const items = Array.isArray(apis) ? apis : [];
    const query = normalizeString(state.apiSearchQuery).toLowerCase();
    if (!query) {
      return items;
    }
    return items.filter((api) => {
      const haystack = [
        api.semantic_name,
        api.method,
        api.path_template,
        api.url,
        api.content_type,
        api.resource_type,
        api.operation_type,
        api.trigger_route,
        api.trigger_action,
      ]
        .filter(Boolean)
        .join(" ")
        .toLowerCase();
      return haystack.includes(query);
    });
  }

  function renderEmptyCardForCurrentTab() {
    if (state.currentTab === "apis" && state.apis.length && state.apiSearchQuery) {
      return (
        '<article class="wb-empty-card"><h3>没有命中的 API</h3><p>试试搜索路径片段、请求方法或 Content-Type，例如 `search`、`GET`、`json`。</p></article>'
      );
    }
    if (state.currentTab === "entities") {
      return '<article class="wb-empty-card"><h3>当前没有实体卡片</h3><p>这个站点暂时还没有沉淀出结构化实体。</p></article>';
    }
    return '<article class="wb-empty-card"><h3>当前没有卡片</h3><p>先执行探索，或者换一个已经有知识沉淀的站点。</p></article>';
  }

  function renderSiteContext() {
    const site = findSite(state.currentSiteId);
    const fallbackName =
      normalizeString(refs.siteNameInput.value) || normalizeString(refs.siteIdInput.value) || "未选择";
    const siteName = site ? site.name || site.site_id || fallbackName : fallbackName;
    const siteURL = site ? site.start_url || "" : normalizeString(refs.siteUrlInput.value);
    const modeLabel =
      state.lastExplore && state.lastExplore.explore_mode
        ? describeExploreMode(state.lastExplore.explore_mode)
        : "未探索";
    const latestExplore =
      state.lastExplore && state.lastExplore.finished_at ? formatDate(state.lastExplore.finished_at) : "暂无";
    const countsLabel =
      state.pages.length || state.apis.length || state.entities.length
        ? state.pages.length +
          " 页面 / " +
          state.apis.length +
          " API / " +
          state.entities.length +
          " 实体"
        : "还没有探索结果";
    const allowedDomains = site && Array.isArray(site.allowed_domains) ? site.allowed_domains.join(", ") : normalizeString(refs.siteDomainsInput.value) || "未设置";
    const sessionName = normalizeString((site && site.session_name) || refs.siteSessionSelect.value) || "未设置";
    const providerName = normalizeString((site && site.provider_id) || refs.siteProviderSelect.value) || "自动选择";
    const runID = state.lastExplore && state.lastExplore.run_id ? state.lastExplore.run_id : "暂无";
    const artifactRoot = state.appMeta && state.appMeta.artifact_root ? state.appMeta.artifact_root : state.artifactRoot || "default";

    refs.currentSiteHeader.textContent = siteName;
    refs.currentSiteStatus.textContent = describeCurrentSiteStatus();

    if (!site && !siteURL) {
      refs.currentSiteOverview.innerHTML =
        '<article class="wb-empty-card"><h3>还没有选中站点</h3><p>先展开“切换站点与接入配置”，选择一个已保存站点，或者填一个地址保存下来。</p></article>';
      return;
    }

    refs.currentSiteOverview.innerHTML =
      '<article class="wb-site-overview-card">' +
      '<div class="wb-site-main">' +
      "<div><h3>" +
      escapeHTML(siteName) +
      "</h3></div>" +
      '<p class="wb-site-url">' +
      escapeHTML(siteURL || "当前还没有保存 Start URL") +
      "</p>" +
      '<div class="wb-site-detail-grid">' +
      siteDetail("探索模式", modeLabel) +
      siteDetail("最近探索", latestExplore) +
      siteDetail("结果", countsLabel) +
      "</div>" +
      '<details class="wb-details wb-details-subtle wb-site-advanced">' +
      "<summary>高级信息</summary>" +
      '<div class="wb-detail-grid">' +
      siteDetail("Allowed Domains", allowedDomains) +
      siteDetail("Session", sessionName) +
      siteDetail("Provider", providerName) +
      (isDeveloperMode() ? siteDetail("Run ID", runID) + siteDetail("Artifact Root", artifactRoot) : "") +
      "</div>" +
      "</details>" +
      "</div>" +
      "</article>";
  }

  function siteDetail(label, value) {
    return (
      '<div class="wb-site-detail-item"><span>' +
      escapeHTML(label) +
      "</span><strong>" +
      escapeHTML(value) +
      "</strong></div>"
    );
  }

  function miniStat(label, value) {
    return (
      '<div class="wb-mini-stat"><span>' +
      escapeHTML(label) +
      "</span><strong>" +
      escapeHTML(value) +
      "</strong></div>"
    );
  }

  function isDeveloperMode() {
    return state.viewMode === "developer";
  }

  function handleFlowEditorInput() {
    state.currentFlowYAML = refs.flowEditor.value || "";
    syncFlowActionButtons();
  }

  function setCurrentFlowYAML(value) {
    state.currentFlowYAML = typeof value === "string" ? value : "";
    refs.flowEditor.value = state.currentFlowYAML || "# 还没有生成 Flow";
    syncFlowActionButtons();
  }

  function currentFlowYAMLValue() {
    return typeof state.currentFlowYAML === "string" ? state.currentFlowYAML : "";
  }

  function originalFlowYAMLValue() {
    return state.lastPlan && typeof state.lastPlan.flow_yaml === "string" ? state.lastPlan.flow_yaml : "";
  }

  function syncFlowActionButtons() {
    const hasFlow = !!normalizeString(currentFlowYAMLValue());
    refs.runFlowButton.disabled = !hasFlow;
    refs.replayFlowButton.disabled = !hasFlow;
    refs.runRepairedFlowButton.disabled = !hasFlow;
  }

  function normalizeViewMode(value) {
    return normalizeString(value) === "developer" ? "developer" : "novice";
  }

  function readStoredViewMode() {
    try {
      return normalizeViewMode(window.localStorage.getItem("tsplay.workbench.viewMode"));
    } catch (_error) {
      return "novice";
    }
  }

  function writeStoredViewMode(mode) {
    try {
      window.localStorage.setItem("tsplay.workbench.viewMode", normalizeViewMode(mode));
    } catch (_error) {
      // Ignore storage failures.
    }
  }

  function totalPageEventCount(pages) {
    return (pages || []).reduce((total, item) => total + (Array.isArray(item.events) ? item.events.length : 0), 0);
  }

  function countWarningEvents(events) {
    return (events || []).reduce((total, item) => {
      const level = normalizeString(item && item.level).toLowerCase();
      return total + (level.includes("warn") ? 1 : 0);
    }, 0);
  }

  function summarizeEventLevels(events) {
    return (events || []).reduce(
      (summary, item) => {
        const level = normalizeString(item && item.level).toLowerCase();
        if (level.includes("error")) {
          summary.error += 1;
        } else if (level.includes("warn")) {
          summary.warning += 1;
        } else {
          summary.info += 1;
        }
        return summary;
      },
      { warning: 0, info: 0, error: 0 }
    );
  }

  function describeCurrentSiteStatus() {
    const hasSite = !!normalizeString(state.currentSiteId);
    const hasKnowledge =
      state.pages.length || state.apis.length || state.entities.length || (state.lastExplore && state.lastExplore.finished_at);
    if (!hasSite) {
      return "未选择";
    }
    if (hasKnowledge) {
      return "已探索";
    }
    return "待探索";
  }

  function describeExploreHeadline() {
    if (!(state.pages.length || state.apis.length || state.entities.length)) {
      return "站点探索已完成";
    }
    return (
      "已完成探索：" +
      state.pages.length +
      " 个页面，" +
      state.apis.length +
      " 个 API，" +
      state.entities.length +
      " 个实体，" +
      totalPageEventCount(state.pages) +
      " 条事件"
    );
  }

  function describeObservationStatus(summary) {
    const technical = normalizeObservationTechnicalDetail(summary && summary.observation_summary);
    const observationErrorCount = Number(summary && summary.observation_error_count ? summary.observation_error_count : summary && summary.observationErrorCount ? summary.observationErrorCount : 0);
    if (!summary) {
      return {
        message: "页面结构观察信息暂不可用。",
        technical: "",
      };
    }
    if (observationErrorCount > 0 || technical) {
      return {
        message: "页面结构观察已降级，已优先保留接口识别和探索结果保存，页面元素识别可能不完整。",
        technical: technical,
      };
    }
    return {
      message: "已完成页面结构观察。",
      technical: "",
    };
  }

  function normalizeObservationTechnicalDetail(value) {
    const text = normalizeString(value);
    if (!text) {
      return "";
    }
    return /observe skipped|playwright rpc|safe mode|fallback|error|failed|降级/i.test(text) ? text : "";
  }

  function summarizeObservationNote(summary) {
    if (!summary) {
      return "";
    }
    return describeObservationStatus(summary).message;
  }

  function summarizePageObservation(summary) {
    if (!summary) {
      return "";
    }
    return describeObservationStatus(summary).message;
  }

  function describePayloadType(contentType, resourceType) {
    const resource = normalizeString(resourceType);
    if (resource && resource !== "unknown") {
      return resource.toUpperCase();
    }
    const value = normalizeString(contentType).toLowerCase();
    if (!value) {
      return "";
    }
    if (value.includes("json")) {
      return "JSON";
    }
    if (value.includes("html")) {
      return "HTML";
    }
    if (value.includes("xml")) {
      return "XML";
    }
    if (value.includes("javascript")) {
      return "JS";
    }
    if (value.includes("form-urlencoded")) {
      return "FORM";
    }
    if (value.includes("multipart")) {
      return "MULTIPART";
    }
    if (value.startsWith("text/")) {
      return "TEXT";
    }
    return value.split(";")[0].toUpperCase();
  }

  function describeBusyStatus(label) {
    switch (label) {
      case "加载 Workbench":
        return "正在加载 Workbench…";
      case "刷新数据":
        return "正在刷新工作台数据…";
      case "探索站点":
        return "正在探索站点…";
      default:
        return "正在" + label + "…";
    }
  }

  function describeSuccessStatus(label) {
    switch (label) {
      case "加载 Workbench":
        return "Workbench 已加载";
      case "刷新数据":
        return "工作台数据已刷新";
      case "加载站点":
      case "切换站点":
        return "当前站点已切换";
      case "保存站点":
        return "站点已保存";
      case "保存 Session":
        return "Session 已保存";
      case "保存 Provider":
        return "Provider 已保存";
      case "探索站点":
        return describeExploreHeadline();
      case "生成 Flow":
        return state.lastPlan && normalizeString(state.lastPlan.flow_yaml) ? "Flow 已生成，可先查看步骤卡片。" : "任务规划已完成。";
      case "执行 Flow":
        return state.lastRun && state.lastRun.ok ? "Flow 已执行完成。" : "Flow 执行结束。";
      case "回放 Flow":
        return state.lastRun && state.lastRun.ok ? "Flow 已回放完成。" : "Flow 回放结束。";
      case "生成 Repair Context":
        return "Repair Context 已生成。";
      case "自动修复 Flow":
        return "Flow 自动修复已完成。";
      case "执行修复版 Flow":
        return state.lastRun && state.lastRun.ok ? "修复版 Flow 已执行完成。" : "修复版 Flow 执行结束。";
      case "复制 Flow":
        return "Flow 已复制。";
      default:
        return label + "完成";
    }
  }

  function metaCard(label, value, tabTarget) {
    const tagName = tabTarget ? "button" : "div";
    const attrs = tabTarget
      ? ' class="wb-meta-card wb-meta-card-action" type="button" data-tab-target="' + escapeAttr(tabTarget) + '"'
      : ' class="wb-meta-card"';
    return (
      "<" +
      tagName +
      attrs +
      "><span>" +
      escapeHTML(label) +
      "</span><strong>" +
      escapeHTML(value) +
      "</strong></" +
      tagName +
      ">"
    );
  }

  function pill(text) {
    return '<span class="wb-pill">' + text + "</span>";
  }

  function describeExploreMode(mode) {
    switch (normalizeString(mode)) {
      case "authorized_dom_api":
        return "已登录后台 DOM/API";
      case "public_html_fallback":
        return "公开页 HTML";
      default:
        return mode || "未知";
    }
  }

  function describePageEventType(type) {
    switch (normalizeString(type)) {
      case "frame_navigated":
        return "页面跳转";
      case "popup":
        return "子窗口";
      case "download":
        return "下载";
      case "console":
        return "Console";
      case "page_error":
        return "页面报错";
      case "websocket":
        return "WebSocket";
      case "websocket_error":
        return "WebSocket 报错";
      case "websocket_closed":
        return "WebSocket 关闭";
      default:
        return type || "事件";
    }
  }

  function formatRiskLabel(risk) {
    const normalized = normalizeString(risk) || "unknown";
    switch (normalized) {
      case "read":
        return "read";
      case "read_download":
        return "read_download";
      case "write_low":
        return "write_low";
      case "write_high":
        return "write_high";
      case "critical":
        return "critical";
      default:
        return normalized;
    }
  }

  function riskPill(risk) {
    const normalized = formatRiskLabel(risk);
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

  function deriveLatestExploreRunID(pages) {
    const items = Array.isArray(pages) ? pages : [];
    for (const item of items) {
      const runID = normalizeString(item && item.explore_run_id);
      if (runID) {
        return runID;
      }
    }
    return "";
  }

  function deriveLatestKnowledgeTimestamp(pages, apis, entities) {
    const candidates = [];
    for (const item of [].concat(pages || [], apis || [], entities || [])) {
      const updatedAt = normalizeString(item && item.updated_at);
      if (updatedAt) {
        candidates.push(updatedAt);
      }
    }
    candidates.sort();
    return candidates.length ? candidates[candidates.length - 1] : "";
  }

  function summarizeExploreLifecycle(pages) {
    const summary = {
      filterRule: "仅保留 xhr/fetch，请求阶段自动排除 html/css/js/image/font/media 等静态资源。",
      observationMode: "页面 load 完成后再做渲染态观察，补抓交互元素、内容块、截图和 DOM 快照。",
      networkRequestCount: 0,
      readableResponseCount: 0,
      networkFailureCount: 0,
      eventCount: 0,
      interactiveElementCount: 0,
      contentElementCount: 0,
      observationErrorCount: 0,
      observationSummary: "",
    };
    (pages || []).forEach((page) => {
      const capture = page && page.capture_summary ? page.capture_summary : null;
      if (!capture) {
        return;
      }
      summary.filterRule = capture.filter_rule || summary.filterRule;
      summary.observationMode = capture.observation_mode || summary.observationMode;
      summary.networkRequestCount += Number(capture.network_request_count || 0);
      summary.readableResponseCount += Number(capture.readable_response_count || 0);
      summary.networkFailureCount += Number(capture.network_failure_count || 0);
      summary.eventCount += Number(capture.event_count || 0);
      summary.interactiveElementCount += Number(capture.interactive_element_count || 0);
      summary.contentElementCount += Number(capture.content_element_count || 0);
      summary.observationErrorCount += Number(capture.observation_error_count || 0);
      if (!summary.observationSummary && normalizeString(capture.observation_summary)) {
        summary.observationSummary = normalizeString(capture.observation_summary);
      }
    });
    return summary;
  }

  function knowledgeSignature() {
    const pagePart = (state.pages || [])
      .map((item) => [item.id || "", item.explore_run_id || "", item.updated_at || ""].join("@"))
      .sort()
      .join("|");
    const apiPart = (state.apis || [])
      .map((item) => [item.id || "", item.updated_at || "", item.status || ""].join("@"))
      .sort()
      .join("|");
    const entityPart = (state.entities || [])
      .map((item) => [item.id || "", item.updated_at || ""].join("@"))
      .sort()
      .join("|");
    return [pagePart, apiPart, entityPart].join("::");
  }

  function sleep(ms) {
    return new Promise((resolve) => {
      window.setTimeout(resolve, ms);
    });
  }

  function appendCacheBuster(path) {
    const value = normalizeString(path);
    if (!value) {
      return value;
    }
    const joiner = value.includes("?") ? "&" : "?";
    return value + joiner + "_ts=" + Date.now();
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

  function firstNonEmpty() {
    for (let index = 0; index < arguments.length; index += 1) {
      const value = normalizeString(arguments[index]);
      if (value) {
        return value;
      }
    }
    return "";
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

  async function loadArtifactJSON(path) {
    const href = artifactHrefForPath(path);
    if (!href) {
      throw new Error("artifact path is not available from the current app context");
    }
    const response = await fetch(appendCacheBuster(href), {
      cache: "no-store",
    });
    if (!response.ok) {
      throw new Error("artifact request failed with status " + response.status);
    }
    return response.json();
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
        cache: "no-store",
      },
      options || {}
    );
    if (requestOptions.body && !requestOptions.headers["Content-Type"]) {
      requestOptions.headers["Content-Type"] = "application/json";
    }
    const method = normalizeString(requestOptions.method || "GET").toUpperCase();
    const requestPath = method === "GET" || method === "HEAD" ? appendCacheBuster(path) : path;
    const response = await fetch(requestPath, requestOptions);
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
