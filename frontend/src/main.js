import './style.css';
import { GetPorts, StartEmission, StopEmission } from '../wailsjs/go/main/App';
import * as runtime from '../wailsjs/runtime/runtime';

const app = document.querySelector('#app');

let state = {
    ports: [],
    selectedPort: '',
    baudRate: 9600,
    sentences: ['GGA', 'RMC'],
    frequency: 1,
    isEmitting: false,
    sentencesSent: 0,
    startTime: null,
    history: []
};

function render() {
    app.innerHTML = `
        <div class="card sidebar">
            <h1>NMEA Gen</h1>
            
            <div class="form-group">
                <label>Port</label>
                <select id="port-select">
                    <option value="">Select a port...</option>
                    ${state.ports.map(p => `<option value="${p}" ${p === state.selectedPort ? 'selected' : ''}>${p}</option>`).join('')}
                </select>
                <button id="refresh-ports" style="margin-top: 10px; padding: 5px; font-size: 0.7rem; background: var(--panel-bg); color: var(--text-dim);">Refresh Ports</button>
            </div>

            <div class="form-group">
                <label>Baud Rate</label>
                <select id="baud-select">
                    ${[4800, 9600, 19200, 38400, 57600, 115200].map(b => `<option value="${b}" ${b === state.baudRate ? 'selected' : ''}>${b}</option>`).join('')}
                </select>
            </div>

            <div class="form-group">
                <label>Sentences</label>
                <div class="checkbox-group">
                    ${['GGA', 'RMC', 'VTG'].map(s => `
                        <div class="checkbox-item ${state.sentences.includes(s) ? 'active' : ''}" data-sentence="${s}">
                            ${s}
                        </div>
                    `).join('')}
                </div>
            </div>

            <div class="form-group">
                <label>Frequency (Hz): <span id="freq-val">${state.frequency}</span></label>
                <input type="range" id="freq-range" min="1" max="10" step="1" value="${state.frequency}" style="width: 100%;">
            </div>

            <button id="toggle-btn" class="${state.isEmitting ? 'stop' : ''}">
                ${state.isEmitting ? 'Stop Emission' : 'Start Emission'}
            </button>

            <div class="stats">
                <div class="stat-box">
                    <span class="stat-value">${state.sentencesSent}</span>
                    <span class="stat-label">Sentences</span>
                </div>
                <div class="stat-box">
                    <span class="stat-value" id="uptime">0s</span>
                    <span class="stat-label">Uptime</span>
                </div>
            </div>
        </div>

        <div class="card main-view">
            <h2>Live Monitor</h2>
            <div id="nmea-history" class="sentence-list">
                ${state.history.map(s => `<div class="sentence-item">${s}</div>`).join('')}
            </div>
        </div>

        <div class="card log-view">
            <h3>System Status</h3>
            <div id="system-log" style="font-size: 0.8rem; color: var(--text-dim); overflow-y: auto;">
                <div>Application Ready.</div>
            </div>
        </div>
    `;

    setupEvents();
}

function setupEvents() {
    document.querySelector('#port-select').onchange = (e) => state.selectedPort = e.target.value;
    document.querySelector('#baud-select').onchange = (e) => state.baudRate = parseInt(e.target.value);
    document.querySelector('#freq-range').oninput = (e) => {
        state.frequency = parseFloat(e.target.value);
        document.querySelector('#freq-val').innerText = state.frequency;
    };
    
    document.querySelectorAll('.checkbox-item').forEach(item => {
        item.onclick = () => {
            const s = item.dataset.sentence;
            if (state.sentences.includes(s)) {
                state.sentences = state.sentences.filter(x => x !== s);
            } else {
                state.sentences.push(s);
            }
            render();
        };
    });

    document.querySelector('#refresh-ports').onclick = async () => {
        log("Refreshing ports...");
        state.ports = await GetPorts() || [];
        render();
    };

    document.querySelector('#toggle-btn').onclick = async () => {
        if (state.isEmitting) {
            await StopEmission();
            state.isEmitting = false;
            state.startTime = null;
            log("Emission stopped.");
        } else {
            if (!state.selectedPort) {
                alert("Please select a port first.");
                return;
            }
            try {
                await StartEmission(state.selectedPort, state.baudRate, state.sentences, state.frequency);
                state.isEmitting = true;
                state.sentencesSent = 0;
                state.startTime = Date.now();
                log(`Started emission on ${state.selectedPort} at ${state.baudRate} baud.`);
            } catch (err) {
                log(`Error: ${err}`, true);
                alert(err);
            }
        }
        render();
    };
}

function log(msg, isError = false) {
    const logDiv = document.querySelector('#system-log');
    if (logDiv) {
        const entry = document.createElement('div');
        entry.style.color = isError ? '#ff4d4d' : 'var(--text-dim)';
        entry.innerText = `[${new Date().toLocaleTimeString()}] ${msg}`;
        logDiv.prepend(entry);
    }
}

// Global Wails Event listeners
runtime.EventsOn("nmea-sentence", (sentence) => {
    state.sentencesSent++;
    state.history.unshift(sentence);
    if (state.history.length > 50) state.history.pop();
    
    const historyDiv = document.querySelector('#nmea-history');
    if (historyDiv) {
        const item = document.createElement('div');
        item.className = 'sentence-item';
        item.innerText = sentence;
        historyDiv.prepend(item);
        if (historyDiv.children.length > 50) historyDiv.lastChild.remove();
    }
    
    document.querySelectorAll('.stat-value')[0].innerText = state.sentencesSent;
});

runtime.EventsOn("log", (msg) => {
    log(msg, true);
});

// Uptime timer
setInterval(() => {
    if (state.isEmitting && state.startTime) {
        const diff = Math.floor((Date.now() - state.startTime) / 1000);
        const uptimeEl = document.querySelector('#uptime');
        if (uptimeEl) uptimeEl.innerText = `${diff}s`;
    }
}, 1000);

// Initial Load
(async () => {
    try {
        state.ports = await GetPorts() || [];
    } catch (e) {
        log("Failed to fetch ports", true);
    }
    render();
})();
