async function loadProgress() {
  const res = await fetch("data/progress.json");
  const json = await res.json();

  renderProgressBar(json.percent, json.done, json.total);
  renderProgressTable(json.genres);
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
    <table>
      <thead>
        <tr>
          <th>Genre</th>
          <th>RAW</th>
          <th>PNG</th>
          <th>ICO</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
  `;

  genres.forEach(genre => {
    html += `
      <tr>
        <td>${genre.genre}</td>
        <td>${genre.raw}</td>
        <td>${genre.png}</td>
        <td>${genre.ico}</td>
        <td>${genre.status}</td>
      </tr>
    `;
  });

  html += `
      </tbody>
    </table>
  `;

  container.innerHTML = html;
}

window.addEventListener("DOMContentLoaded", loadProgress);
