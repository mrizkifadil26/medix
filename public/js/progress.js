async function loadProgress() {
  const res = await fetch("data/progress.json");
  const json = await res.json();

  renderProgressBar(json.percent, json.done, json.total);
  renderProgressTable(json.genres);
  renderProgressCards(json.genres);
  renderProgressLegend();
}

function renderProgressBar(percent, done, total) {
  const container = document.getElementById("progress-bar-container");
  container.innerHTML = `
    <div class="progress-container">
      <strong>Progress:</strong>
      <div class="progress-bar">
        <div class="progress-fill" style="width: ${percent}%;"></div>
        <div class="progress-text">${percent}% (${done}/${total})</div>
      </div>
    </div>
  `;
}

function renderProgressTable(genres) {
  const container = document.getElementById("progress-table");

  let html = `
    <table class="progress-table">
      <thead>
        <tr>
          <th>Genre</th>
          <th>Progress</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
  `;

  genres.forEach(g => {
    const raw = g.raw || 1;
    const ico = g.ico || 0;
    const png = g.png || 0;
    const pngOnly = png - ico;
    const rawOnly = raw - png;

    const icoWidth = (ico / raw) * 100;
    const pngWidth = (pngOnly / raw) * 100;
    const rawWidth = (rawOnly / raw) * 100;

    // ✅ Correct percentage is based on combined contribution
    const pngPercent = (png / raw) * 100;
    const percent = Math.min(Math.round((ico / raw) * 100), 100);
    console.log(`📂 Genre: ${g.genre || "Unknown"}
      Raw Total     : ${raw}
      ICO Count     : ${ico}
      PNG Count     : ${png}
      PNG Only      : ${pngOnly}
      Raw Only      : ${rawOnly}
      ICO Width     : ${icoWidth.toFixed(2)}%
      PNG Width     : ${pngWidth.toFixed(2)}%
      Raw Width     : ${rawWidth.toFixed(2)}%
      Progress      : ${percent}%
    `);

    let labelClass = "bar-label";
    if (percent < 20 && pngPercent < 40) {
      labelClass += " label-light label-right";
    } else if (percent < 40 && pngPercent < 40) {
      labelClass += " label-light";
    }

    html += `
      <tr>
        <td data-label="Genre">
          <span class="genre-label">
            ${g.icon || "🎬"} ${g.genre}
          </span>
        </td>

        <td data-label="Progress">
          <div class="bar-wrapper">
            <div class="bar-container">
              <div class="bar-ico" style="width: ${icoWidth}%"></div>
              <div class="bar-png" style="width: ${pngWidth}%"></div>
              <div class="bar-raw" style="width: ${rawWidth}%"></div>
              <div class="${labelClass}">${percent}%</div>
            </div>
            <div class="bar-info">
              <span class="bar-count ico">${ico}</span> /
              <span class="bar-count png">${png}</span> /
              <span class="bar-count raw">${raw}</span>
            </div>
          </div>
        </td>
        <td data-label="Status">${g.status}</td>
      </tr>
    `;
  });

  html += `</tbody></table>`;
  container.innerHTML = html;
}

function renderProgressCards(genres) {
  const container = document.getElementById("progress-cards");
  container.innerHTML = "";

  genres.forEach(g => {
    const raw = g.raw || 1;
    const ico = g.ico || 0;
    const png = g.png || 0;
    const pngOnly = png - ico;
    const rawOnly = raw - png;

    const icoWidth = (ico / raw) * 100;
    const pngWidth = (pngOnly / raw) * 100;
    const rawWidth = (rawOnly / raw) * 100;

    const percent = Math.round((ico / raw) * 100);

    const card = document.createElement("div");
    card.className = "genre-card";
    card.innerHTML = `
      <div class="genre-header">
        <div class="genre-name">${g.icon || "🎬"} ${g.genre}</div>
        <div class="genre-status">${g.status}</div>
      </div>
      <div class="bar-container">
        <div class="bar-ico" style="width: ${icoWidth}%;"></div>
        <div class="bar-png" style="width: ${pngWidth}%;"></div>
        <div class="bar-raw" style="width: ${rawWidth}%;"></div>
        <div class="bar-label">${Math.round(percent)}%</div>
      </div>

      <div class="bar-info">
        <span class="bar-count ico">${ico}</span> /
        <span class="bar-count png">${png}</span> /
        <span class="bar-count raw">${raw}</span>
      </div>
    `;
    container.appendChild(card);
  });
}

function renderProgressLegend() {
  const legendContainer = document.getElementById("progress-legend");
  legendContainer.innerHTML = `
    <div class="legend">
      <strong>Legend:</strong>
      <div><span class="box bar-ico"></span> ICO</div>
      <div><span class="box bar-png"></span> PNG</div>
      <div><span class="box bar-raw"></span> RAW</div>
    </div>
  `;
}

function createLabel(percentage, text) {
  const label = document.createElement("div");
  label.classList.add("bar-label");

  if (percentage < 20) {
    label.classList.add("label-light", "label-right");
  } else if (percentage < 40) {
    label.classList.add("label-light");
  }

  label.textContent = text;
  return label;
}

window.addEventListener("DOMContentLoaded", loadProgress);