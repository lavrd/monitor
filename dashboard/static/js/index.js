// charts update interval
let interval = null;

// chart update interval time
const updateChartInterval = 3000,
  // api address
  api = 'http://localhost:4222/api/metrics/';

// show charts
const show = () => {
  const ids = document.getElementById('containerName');
  if (ids.value === '') {
    ids.value = 'all';
  }
  stop();
  interval = showCharts(ids.value);
  ids.value = '';
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
    fetch(api + ids).then(response => {
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
    grid: {x: {show: !0,}, y: {show: !0,}}
  });
};
