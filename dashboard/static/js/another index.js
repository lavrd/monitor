// charts update interval
let interval = null;

// Array.from({length: 6}, () => Math.floor(Math.random() * 100))
function getRandomColor() {
  var letters = '0123456789ABCDEF';
  var color = '#';
  for (var i = 0; i < 6; i++) {
    color += letters[Math.floor(Math.random() * 16)];
  }
  return color;
}

const randomdata = () => {
  return Array.from({length: 25}, () => Math.floor(Math.random() * 100));
};


window.onload = () => {

  var ctx = document.getElementById("myChart").getContext('2d');
  var ctx2 = document.getElementById("myChart2").getContext('2d');
  // var ctx3 = document.getElementById("myChart3").getContext('2d');
  const array1 = randomdata();
  const color1 = getRandomColor();
  const config = {
    type: 'line',
    data: {
      labels: randomdata(),
      datasets: [{
        label: 'gnatsd',
        backgroundColor: color1,
        borderColor: color1,
        data: array1,
        fill: false
      }, {
        label: 'gnatsd',
        backgroundColor: color1,
        borderColor: color1,
        data: array1,
        fill: false
      }]
    },
    options: {
      elements: {
        line: {
          // tension: 0 // disables bezier curves
        }
      },
      animation: {
        // duration: 0 // general animation time
      },
      responsiveAnimationDuration: 0, // animation duration after a resize

      responsive: true,
      title: {
        display: true,
        text: 'Chart.js Line Chart'
      },
      tooltips: {
        mode: 'index',
        intersect: false
      },
      hover: {
        mode: 'nearest',
        animationDuration: 0, // duration of animations when hovering an item
        intersect: true
      },
      scales: {
        xAxes: [{
          display: true,
          scaleLabel: {
            display: true,
            labelString: 'Month'
          }
        }],
        yAxes: [{
          display: true,
          scaleLabel: {
            display: true,
            labelString: 'Value'
          }
        }]
      }
    }
  };
  const chart = new Chart(ctx, config);
  const chart2 = new Chart(ctx2, config);
  // const chart3 = new Chart(ctx3, config);

  setInterval(() => {
    config.data.datasets.forEach(function (dataset) {
      dataset.data = randomdata();
    });

    chart.update();
    chart2.update();
    // chart3.update();
  }, 1000);

};

// const socket = new WebSocket("ws://localhost:2000/metrics");

// chart update interval time
const updateChartInterval = 3000,
  // api address
  api = 'http://localhost:4222/api/metrics';

// socket.onopen = function () {
//   alert("Соединение установлено.");
//   console.log(1);
//   socket.send("preved");
// };
//
// socket.onclose = function (event) {
//   if (event.wasClean) {
//     alert('Соединение закрыто чисто');
//   } else {
//     alert('Обрыв соединения'); // например, "убит" процесс сервера
//   }
//   alert('Код: ' + event.code + ' причина: ' + event.reason);
// };
//
// socket.onmessage = function (event) {
//   console.log(event.data.metrics);
// };
//
// socket.onerror = function (error) {
//   alert("Ошибка " + error.message);
// };
// show charts
const show = () => {
  // const ids = document.getElementById('containerName');
  // if (ids.value === '') {
  //   ids.value = 'all';
  // }
  // stop();
  // interval = showCharts(ids.value);
  // ids.value = '';

  socket.send("kek");
};

// clear all charts
const clearCharts = () => {
  const elements = document.getElementsByClassName("temp");
  for (let i = elements.length - 1; i >= 0; i--) {
    if (elements[i] && elements[i].parentElement) {
      elements[i].parentElement.removeChild(elements[i]);
    }
  }
};

// stop / clear charts
const stop = () => {
  changeNotification();
  clearInterval(interval);
  clearCharts();
};

// create chart div
const createChartDiv = (parent, name) => {
  const h2 = document.createElement('h2');
  h2.innerText = name;
  h2.setAttribute('class', 'temp title');
  h2.setAttribute('id', 'h2' + name);
  parent.appendChild(h2);

  const div = document.createElement('div');
  div.setAttribute('id', name);
  div.setAttribute('class', 'temp');
  parent.appendChild(div);
};

// remove chart div
const removeChartDiv = (id) => {
  const div = document.getElementById(id);
  const h2 = document.getElementById('h2' + id);
  if (div && div.parentElement) {
    div.parentElement.removeChild(div);
    h2.parentElement.removeChild(h2);
  }
};

// show and update charts
const showCharts = (ids) => {
  let chart = new Map(),
    cpu = new Map(),
    mem = new Map(),
    time = new Map();

  return setInterval(() => {
    fetch(api + ids, {}).then(response => {
      return response.json();
    }).then(data => {
      // if there is stopped containers
      if (data.stopped) {
        for (let i in data.stopped) {
          if (ids.includes(data.stopped[i]) || ids === 'all') {
            // remove container div
            removeChartDiv(data.stopped[i]);

            // todo use for each
            // clear container map
            cpu.delete(data.stopped[i]);
            mem.delete(data.stopped[i]);
            time.delete(data.stopped[i]);
            chart.delete(data.stopped[i]);
          }
        }
      }

      if (data.message) {
        throw data.message;
      }

      // update charts
      for (let i in data.metrics) {
        const id = data.metrics[i].Name;

        // if container chart already exists
        if (chart.has(id)) {
          // update data
          cpu.set(id, setData(cpu.get(id), 'cpu', data.metrics[i].CPUPercentage));
          mem.set(id, setData(mem.get(id), 'mem', data.metrics[i].MemoryPercentage));
          time.set(id, setData(time.get(id), 'time', new Date()));

          // update chart
          chart.get(id).load({
            columns: [time.get(id), cpu.get(id), mem.get(id)]
          });
        } else {
          // if container chart not exists
          // create chart div
          createChartDiv(document.getElementById('chart'), id);

          // init arrays
          cpu.set(id, setData(['cpu'], 'cpu', 0));
          mem.set(id, setData(['mem'], 'mem', 0));
          time.set(id, setData(['time'], 'time', new Date()));

          // show chart
          chart.set(id, createChart(
            id,
            time.get(id),
            cpu.get(id),
            mem.get(id)
          ));
        }
      }

      changeNotification();
    }).catch(error => {
      cpu.clear();
      mem.clear();
      time.clear();
      chart.clear();

      changeNotification(error);
    });
  }, updateChartInterval);
};

// update cpu, mem, time data array
const setData = (data, type, value) => {
  if (data.length === 10) {
    data.shift();
    data.shift();
    data.unshift(type);
  }
  data.push(value);

  return data;
};

// change notification status
const changeNotification = (error) => {
  const alertErrorText = document.getElementById('alert');

  if (!error) {
    alertErrorText.setAttribute('class', 'text-hide');
    return;
  }

  alertErrorText.setAttribute('class', 'alert alert-danger');
  alertErrorText.innerText = error;
};

// create new chart
const createChart = (id, time, cpu, mem) => {
  return c3.generate({
    bindto: '#' + id,
    data: {x: 'time', columns: [time, cpu, mem], type: 'spline'},
    axis: {x: {type: 'timeseries', tick: {format: '%H:%M:%S'}}, y: {tick: {format: d3.format(',.2f')}, label: '%'}},
    grid: {x: {show: !0}, y: {show: !0}}
  });
};
