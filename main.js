// var Chart = require('chart.js');
//import { Chart } from 'chart.js';
//import zoomPlugin from 'chartjs-plugin-zoom';
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
function hourlyPoints(config) {
  const ys = numbers(config);
  const start = new Date().valueOf();
  return ys.map((y, i) => ({x: datefns.addHours(start, i), y}));
}
const rand255 = () => Math.round(Math.random() * 255);
function randomColor(alpha) {
  return 'rgba(' + rand255() + ',' + rand255() + ',' + rand255() + ',' + (alpha || '.3') + ')';
}

// Parsable with datefns.parse(basicdates[n], 'yyyy-MM-dd HH:mm', new Date())
var basicdates = [
    '2021-09-25 14:01'
];


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

const NUMBER_CFG = {count: 500, min: 0, max: 1000};
const data = {
  datasets: [{
    label: 'My First dataset',
    borderColor: randomColor(0.4),
    backgroundColor: randomColor(0.1),
    pointBorderColor: randomColor(0.7),
    pointBackgroundColor: randomColor(0.5),
    pointBorderWidth: 1,
    data: hourlyPoints(NUMBER_CFG),
  }]
};
const scales = {
  x: {
    position: 'bottom',
    type: 'time',
    ticks: {
      autoSkip: true,
      autoSkipPadding: 50,
      maxRotation: 0
    },
    time: {
      displayFormats: {
        hour: 'HH:mm',
        minute: 'HH:mm',
        second: 'HH:mm:ss'
      }
    }
  },
  y: {
    position: 'right',
    ticks: {
      callback: (val, index, ticks) => index === 0 || index === ticks.length - 1 ? null : val,
    },
    grid: {
      borderColor: randomColor(1),
      color: 'rgba( 0, 0, 0, 0.1)',
    },
    title: {
      display: true,
      text: (ctx) => ctx.scale.axis + ' axis',
    }
  },
};
const config = {
  type: 'scatter',
  data: data,
  options: {
    scales: scales,
    plugins: {
      zoom: zoomOptions,
      title: {
        display: true,
        position: 'bottom',
        text: (ctx) => 'Zoom: ' + zoomStatus() + ', Pan: ' + panStatus()
      }
    },
  }
};

var ctx = document.getElementById('myChart').getContext('2d');
var myChart = new Chart(ctx, config);
/*
    {
    type: 'bar',
    data: {
        labels: ['Red', 'Blue', 'Yellow', 'Green', 'Purple', 'Orange'],
        datasets: [{
            label: '# of Votes',
            data: [12, 19, 3, 5, 2, 3],
            backgroundColor: [
                'rgba(255, 99, 132, 0.2)',
                'rgba(54, 162, 235, 0.2)',
                'rgba(255, 206, 86, 0.2)',
                'rgba(75, 192, 192, 0.2)',
                'rgba(153, 102, 255, 0.2)',
                'rgba(255, 159, 64, 0.2)'
            ],
            borderColor: [
                'rgba(255, 99, 132, 1)',
                'rgba(54, 162, 235, 1)',
                'rgba(255, 206, 86, 1)',
                'rgba(75, 192, 192, 1)',
                'rgba(153, 102, 255, 1)',
                'rgba(255, 159, 64, 1)'
            ],
            borderWidth: 1
        }]
    },
    options: {
        scales: {
            y: {
                beginAtZero: true
            }
        }
    }
});
*/
