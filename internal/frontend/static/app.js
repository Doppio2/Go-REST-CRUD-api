const state = {
  equipment: [],
  experiments: [],
  experimentEquipment: new Map(),
};

const elements = {
  equipmentForm: document.getElementById("equipment-form"),
  experimentForm: document.getElementById("experiment-form"),
  equipmentList: document.getElementById("equipment-list"),
  experimentList: document.getElementById("experiment-list"),
  equipmentCount: document.getElementById("equipment-count"),
  experimentCount: document.getElementById("experiment-count"),
  message: document.getElementById("app-message"),
  refreshButton: document.getElementById("refresh-all"),
  themeToggle: document.getElementById("theme-toggle"),
  equipmentTemplate: document.getElementById("equipment-card-template"),
  experimentTemplate: document.getElementById("experiment-card-template"),
};

const THEME_STORAGE_KEY = "lab-inventory-theme";

async function apiRequest(path, options = {}) {
  const response = await fetch(path, {
    headers: {
      "Content-Type": "application/json",
      ...(options.headers || {}),
    },
    ...options,
  });

  if (response.status === 204) {
    return null;
  }

  const contentType = response.headers.get("Content-Type") || "";
  const isJSON = contentType.includes("application/json");
  const payload = isJSON ? await response.json() : null;

  if (!response.ok) {
    const message = payload?.error?.message || "Request failed";
    throw new Error(message);
  }

  return payload?.data ?? null;
}

function showMessage(message, type = "success") {
  elements.message.textContent = message;
  elements.message.className = `message ${type}`;
}

function clearMessage() {
  elements.message.textContent = "";
  elements.message.className = "message hidden";
}

function applyTheme(theme) {
  document.body.dataset.theme = theme;
  elements.themeToggle.textContent = theme === "dark" ? "Light Theme" : "Dark Theme";
  localStorage.setItem(THEME_STORAGE_KEY, theme);
}

function initializeTheme() {
  const savedTheme = localStorage.getItem(THEME_STORAGE_KEY);
  const preferredDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
  applyTheme(savedTheme || (preferredDark ? "dark" : "light"));
}

function toggleTheme() {
  const current = document.body.dataset.theme === "dark" ? "dark" : "light";
  applyTheme(current === "dark" ? "light" : "dark");
}

function formatDate(value) {
  if (!value) {
    return "No timestamp";
  }

  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return date.toLocaleString();
}

async function loadEquipment() {
  state.equipment = (await apiRequest("/equipment")) || [];
}

async function loadExperiments() {
  state.experiments = (await apiRequest("/experiments")) || [];
}

async function loadExperimentEquipment() {
  const entries = await Promise.all(
    state.experiments.map(async (experiment) => {
      const equipment = (await apiRequest(`/experiments/${experiment.id}/equipment`)) || [];
      return [experiment.id, equipment];
    }),
  );

  state.experimentEquipment = new Map(entries);
}

function renderEquipment() {
  elements.equipmentList.innerHTML = "";
  elements.equipmentCount.textContent = `${state.equipment.length} item${state.equipment.length === 1 ? "" : "s"}`;

  if (state.equipment.length === 0) {
    elements.equipmentList.className = "card-grid empty-state";
    elements.equipmentList.textContent = "No equipment yet.";
    return;
  }

  elements.equipmentList.className = "card-grid";

  for (const item of state.equipment) {
    const node = elements.equipmentTemplate.content.firstElementChild.cloneNode(true);
    node.querySelector(".card-title").textContent = item.name;
    node.querySelector(".card-copy").textContent = item.description || "No description";
    node.querySelector(".card-meta").textContent = `Created: ${formatDate(item.creation_date)}`;

    node.querySelector(".card-edit").addEventListener("click", () => editEquipment(item));
    node.querySelector(".card-delete").addEventListener("click", () => deleteEquipment(item.id));

    elements.equipmentList.appendChild(node);
  }
}

function renderExperimentEquipment(experiment, container) {
  const linked = state.experimentEquipment.get(experiment.id) || [];
  container.innerHTML = "";

  if (linked.length === 0) {
    container.className = "linked-equipment-list empty-inline";
    container.textContent = "No equipment assigned.";
    return;
  }

  container.className = "linked-equipment-list";

  for (const item of linked) {
    const chip = document.createElement("div");
    chip.className = "link-chip";
    chip.innerHTML = `
      <span>${item.name}</span>
      <button class="chip-button" type="button">Remove</button>
    `;
    chip.querySelector("button").addEventListener("click", () => unlinkEquipment(experiment.id, item.id));
    container.appendChild(chip);
  }
}

function buildEquipmentOptions(select, experimentId) {
  const assigned = new Set((state.experimentEquipment.get(experimentId) || []).map((item) => item.id));
  const available = state.equipment.filter((item) => !assigned.has(item.id));

  select.innerHTML = "";

  if (available.length === 0) {
    const option = document.createElement("option");
    option.value = "";
    option.textContent = "No available equipment";
    select.appendChild(option);
    select.disabled = true;
    return;
  }

  select.disabled = false;

  for (const item of available) {
    const option = document.createElement("option");
    option.value = String(item.id);
    option.textContent = item.name;
    select.appendChild(option);
  }
}

function renderExperiments() {
  elements.experimentList.innerHTML = "";
  elements.experimentCount.textContent = `${state.experiments.length} item${state.experiments.length === 1 ? "" : "s"}`;

  if (state.experiments.length === 0) {
    elements.experimentList.className = "card-grid empty-state";
    elements.experimentList.textContent = "No experiments yet.";
    return;
  }

  elements.experimentList.className = "card-grid";

  for (const experiment of state.experiments) {
    const node = elements.experimentTemplate.content.firstElementChild.cloneNode(true);
    const select = node.querySelector(".equipment-select");
    const linkedList = node.querySelector(".linked-equipment-list");
    const form = node.querySelector(".link-form");
    const exportLink = node.querySelector(".experiment-export");

    node.querySelector(".card-title").textContent = experiment.name;
    node.querySelector(".card-copy").textContent = experiment.description || "No description";
    node.querySelector(".card-meta").textContent = `Created: ${formatDate(experiment.creation_date)}`;
    exportLink.href = `/experiments/${experiment.id}/equipment?format=csv`;

    buildEquipmentOptions(select, experiment.id);
    renderExperimentEquipment(experiment, linkedList);

    form.addEventListener("submit", async (event) => {
      event.preventDefault();
      if (!select.value) {
        showMessage("Choose equipment before attaching.", "error");
        return;
      }

      await attachEquipment(experiment.id, Number(select.value));
    });

    node.querySelector(".card-edit").addEventListener("click", () => editExperiment(experiment));
    node.querySelector(".card-delete").addEventListener("click", () => deleteExperiment(experiment.id));

    elements.experimentList.appendChild(node);
  }
}

async function refreshAll(showStatus = false) {
  try {
    clearMessage();
    await loadEquipment();
    await loadExperiments();
    await loadExperimentEquipment();
    renderEquipment();
    renderExperiments();

    if (showStatus) {
      showMessage("Data refreshed.");
    }
  } catch (error) {
    showMessage(error.message, "error");
  }
}

async function createEquipment(event) {
  event.preventDefault();

  const formData = new FormData(elements.equipmentForm);
  const payload = {
    name: String(formData.get("name") || ""),
    description: String(formData.get("description") || ""),
  };

  try {
    await apiRequest("/equipment", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    elements.equipmentForm.reset();
    showMessage("Equipment created.");
    await refreshAll();
  } catch (error) {
    showMessage(error.message, "error");
  }
}

async function createExperiment(event) {
  event.preventDefault();

  const formData = new FormData(elements.experimentForm);
  const payload = {
    name: String(formData.get("name") || ""),
    description: String(formData.get("description") || ""),
  };

  try {
    await apiRequest("/experiments", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    elements.experimentForm.reset();
    showMessage("Experiment created.");
    await refreshAll();
  } catch (error) {
    showMessage(error.message, "error");
  }
}

async function deleteEquipment(id) {
  if (!window.confirm("Delete this equipment item?")) {
    return;
  }

  try {
    await apiRequest(`/equipment/${id}`, { method: "DELETE" });
    showMessage("Equipment deleted.");
    await refreshAll();
  } catch (error) {
    showMessage(error.message, "error");
  }
}

async function deleteExperiment(id) {
  if (!window.confirm("Delete this experiment?")) {
    return;
  }

  try {
    await apiRequest(`/experiments/${id}`, { method: "DELETE" });
    showMessage("Experiment deleted.");
    await refreshAll();
  } catch (error) {
    showMessage(error.message, "error");
  }
}

async function editEquipment(item) {
  const name = window.prompt("Equipment name", item.name);
  if (name === null) {
    return;
  }

  const description = window.prompt("Equipment description", item.description || "");
  if (description === null) {
    return;
  }

  try {
    await apiRequest(`/equipment/${item.id}`, {
      method: "PUT",
      body: JSON.stringify({ name, description }),
    });
    showMessage("Equipment updated.");
    await refreshAll();
  } catch (error) {
    showMessage(error.message, "error");
  }
}

async function editExperiment(item) {
  const name = window.prompt("Experiment name", item.name);
  if (name === null) {
    return;
  }

  const description = window.prompt("Experiment description", item.description || "");
  if (description === null) {
    return;
  }

  try {
    await apiRequest(`/experiments/${item.id}`, {
      method: "PUT",
      body: JSON.stringify({ name, description }),
    });
    showMessage("Experiment updated.");
    await refreshAll();
  } catch (error) {
    showMessage(error.message, "error");
  }
}

async function attachEquipment(experimentId, equipmentId) {
  try {
    await apiRequest(`/experiments/${experimentId}/equipment`, {
      method: "POST",
      body: JSON.stringify({ equipment_id: equipmentId }),
    });
    showMessage("Equipment attached to experiment.");
    await refreshAll();
  } catch (error) {
    showMessage(error.message, "error");
  }
}

async function unlinkEquipment(experimentId, equipmentId) {
  try {
    await apiRequest(`/experiments/${experimentId}/equipment/${equipmentId}`, {
      method: "DELETE",
    });
    showMessage("Equipment unlinked from experiment.");
    await refreshAll();
  } catch (error) {
    showMessage(error.message, "error");
  }
}

elements.equipmentForm.addEventListener("submit", createEquipment);
elements.experimentForm.addEventListener("submit", createExperiment);
elements.refreshButton.addEventListener("click", () => refreshAll(true));
elements.themeToggle.addEventListener("click", toggleTheme);

initializeTheme();
refreshAll();
