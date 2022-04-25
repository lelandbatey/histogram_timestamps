var Chart = require('chart.js');
require('chartjs-plugin-zoom');
var datefns = require('date-fns')
require('chartjs-adapter-date-fns');

// Ripped from the Chartjs source
function valueOrDefault(value, defaultValue) {
  return typeof value === 'undefined' ? defaultValue : value;
}
let _seed = Date.now();
function rand(min, max) {
  min = valueOrDefault(min, 0);
  max = valueOrDefault(max, 0);
  _seed = (_seed * 9301 + 49297) % 233280;
  return min + (_seed / 233280) * (max - min);
}
function numbers(config) {
  const cfg = config || {};
  const min = valueOrDefault(cfg.min, 0);
  const max = valueOrDefault(cfg.max, 100);
  const from = valueOrDefault(cfg.from, []);
  const count = valueOrDefault(cfg.count, 8);
  const decimals = valueOrDefault(cfg.decimals, 8);
  const continuity = valueOrDefault(cfg.continuity, 1);
  const dfactor = Math.pow(10, decimals) || 0;
  const data = [];
  let i, value;

  for (i = 0; i < count; ++i) {
    value = (from[i] || 0) + rand(min, max);
    if (rand() <= continuity) {
      data.push(Math.round(dfactor * value) / dfactor);
    } else {
      data.push(null);
    }
  }

  return data;
}
const rand255 = () => Math.round(Math.random() * 255);


const zoomOptions = {
  pan: {
    enabled: true,
    modifierKey: 'ctrl',
  },
  zoom: {
    drag: {
      enabled: true
    },
    mode: 'xy',
  },
};
// </block>

const panStatus = () => zoomOptions.pan.enabled ? 'enabled' : 'disabled';
const zoomStatus = () => zoomOptions.zoom.drag.enabled ? 'enabled' : 'disabled';


const LABEL_LOCALTZ = 'Timeseries #1 - Local time zone ('+Intl.DateTimeFormat().resolvedOptions().timeZone+')';
const LABEL_UTC = 'Timeseries #1 - UTC';

const LINE_COLOR = 'rgb(54, 162, 235)';

const data = {
    datasets: [
        {
            label: LABEL_LOCALTZ,
            data: CONTEXT.data,
            borderColor: LINE_COLOR,
        }
    ],
};

const config = {
    type: 'line',
    data: data,
    options: {
        parsing: true,
        responsive: true,
        maintainAspectRatio: false,
        animation: {
            // Massively speed up all the default animations
            duration: 200,
        },
        scales: {
            x: {
                type: 'timeseries',
                time: {unit: CONTEXT.unit},
            },
        },
        plugins: {
            zoom: zoomOptions,
            title: {
                display: true,
                position: 'bottom',
                text: (ctx) => 'Zoom: (click and drag)' + zoomStatus() + ', Pan (ctrl + click and drag): ' + panStatus()
            }
        },
    },
};

const ctx = document.getElementById('myChart').getContext('2d');
const myChart = new Chart(ctx, config);
function convertDateToUTC(date_) {
    return new Date(date_.getUTCFullYear(), date_.getUTCMonth(), date_.getUTCDate(), date_.getUTCHours(), date_.getUTCMinutes(), date_.getUTCSeconds());
}
const actions = [
    {
        name: "Set TZ to local timezone",
        handler(chart) {
            let exp_label = LABEL_LOCALTZ;
            if (chart.data.datasets[0].label == exp_label) {
                return;
            }
            chart.data.datasets[0] = {
                label: exp_label,
                data: CONTEXT.data,
                borderColor: LINE_COLOR,
            };
            chart.update();
        },
    },
    {
        name: "Set TZ to UTC",
        handler(chart) {
            let exp_label = LABEL_UTC;
            if (chart.data.datasets[0].label == exp_label) {
                return;
            }
            var nd = [];
            for (var i = 0; i < CONTEXT.data.length; i++) {
                nd.push({
                    x: convertDateToUTC(new Date(CONTEXT.data[i].x)).getTime(),
                    y: CONTEXT.data[i].y,
                });
            }
            console.log(nd);
            chart.data.datasets[0] = {
                label: exp_label,
                data: nd,
                borderColor: LINE_COLOR,
            };
            chart.update();
        },
    },
    {
        name: 'Reset zoom',
        handler(chart) {
            chart.resetZoom();
        },
    }
];

actions.forEach((a, i) => {
  let button = document.createElement("button");
  button.id = "button"+i;
  button.innerText = a.name;
  button.onclick = () => a.handler(myChart);
  document.querySelector(".buttons").appendChild(button);
});
