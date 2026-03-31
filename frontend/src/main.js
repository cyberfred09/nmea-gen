import './style.css';
import { GetPorts, StartEmission, StopEmission } from '../wailsjs/go/main/App';
import * as runtime from '../wailsjs/runtime/runtime';

const app = document.querySelector('#app');

const CATEGORIES = {
    GNSS: ['GGA', 'RMC', 'GLL', 'GSA', 'GSV'],
    Instruments: ['MWV', 'DBT', 'DPT', 'VHW', 'HDM', 'HDT', 'MTW'],
    Autopilot: ['APB', 'BWC', 'BOD', 'XTE'],
    AIS: ['AIS']
};

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
                    <option value="">Select port...</option>
                    ${state.ports.map(p => `<option value="${p}" ${p === state.selectedPort ? 'selected' : ''}>${p}</option>`).join('')}
                </select>
                <button id="refresh-ports" style="margin-top: 5px; padding: 4px; font-size: 0.7rem;">Refresh Ports</button>
            </div>

            <div class="form-group">
                <label>Baud Rate</label>
                <select id="baud-select">
                    ${[4800, 9600, 19200, 38400, 57600, 115200].map(b => `<option value="${b}" ${b === state.baudRate ? 'selected' : ''}>${b}</option>`).join('')}
                </select>
            </div>

            <div class="form-group">
                <label>Frequency (Hz): <span id="freq-val">${state.frequency}</span></label>
                <input type="range" id="freq-range" min="1" max="10" step="1" value="${state.frequency}">
            </div>

            ${Object.entries(CATEGORIES).map(([cat, list]) => `
                <h3>${cat}</h3>
                <div class="checkbox-grid">
                    ${list.map(s => `
                        <div class="checkbox-item ${state.sentences.includes(s) ? 'active' : ''}" data-sentence="${s}">
                            ${s}
                        </div>
                    `).join('')}
                </div>
            `).join('')}

            <button id="toggle-btn" class="${state.isEmitting ? 'stop' : ''}">
                ${state.isEmitting ? 'Stop Emission' : 'Start Emission'}
            </button>

            <div class="stats">
                <div class="stat-box"><span class="stat-value">${state.sentencesSent}</span><span class="stat-label">Sentences</span></div>
                <div class="stat-box"><span class="stat-value" id="uptime">0s</span><span class="stat-label">Uptime</span></div>
            </div>
        </div>

        <div class="card main-view">
            <div id="nmea-history" class="sentence-list">
                ${state.history.map(s => `<div class="sentence-item">${s}</div>`).join('')}
            </div>
        </div>

        <div class="card log-view">
            <div id="system-log" style="font-size: 0.75rem; color: var(--text-dim); overflow-y: auto;">
                <div>Ready.</div>
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
        state.ports = await GetPorts() || [];
        render();
    };

    document.querySelector('#toggle-btn').onclick = async () => {
        if (state.isEmitting) {
            await StopEmission();
            state.isEmitting = false;
            state.startTime = null;
        } else {
            if (!state.selectedPort) return alert("Select a port!");
            try {
                await StartEmission(state.selectedPort, state.baudRate, state.sentences, state.frequency);
                state.isEmitting = true;
                state.sentencesSent = 0;
                state.startTime = Date.now();
            } catch (err) {
                alert(err);
            }
        }
        render();
    };
}

runtime.EventsOn("nmea-sentence", (sentence) => {
    state.sentencesSent++;
    state.history.unshift(sentence);
    if (state.history.length > 100) state.history.pop();
    
    const historyDiv = document.querySelector('#nmea-history');
    if (historyDiv) {
        const item = document.createElement('div');
        item.className = 'sentence-item';
        item.innerText = sentence;
        historyDiv.prepend(item);
        if (historyDiv.children.length > 100) historyDiv.lastChild.remove();
    }
    document.querySelectorAll('.stat-value')[0].innerText = state.sentencesSent;
});

setInterval(() => {
    if (state.isEmitting && state.startTime) {
        const diff = Math.floor((Date.now() - state.startTime) / 1000);
        const el = document.querySelector('#uptime');
        if (el) el.innerText = `${diff}s`;
    }
}, 1000);

(async () => {
    state.ports = await GetPorts() || [];
    render();
})();
