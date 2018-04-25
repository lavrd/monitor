let socket = null;

const charts = new Map();

const elStatus = document.getElementById('status');
const elContainersName = document.getElementById('name');
const elAlert = document.getElementById('alert');
const elRoot = document.getElementById('root');

const STATUS = {
  SUCCESS: 'SUCCESS',
  ERROR: 'ERROR',
  CLOSE: 'CLOSE'
};

const setAlert = (text) => {
  const hideClass = 'text-hide';

  if (elAlert.classList.contains(hideClass) && !!text) elAlert.classList.remove(hideClass);
  if (!text) elAlert.classList.add(hideClass);

  elAlert.innerText = text;
};

const setStatus = (status) => {
  elStatus.className = '';

  const base = 'badge badge-';

  switch (status) {
    case STATUS.SUCCESS:
      elStatus.className = `${base}success`;
      break;
    case STATUS.ERROR:
      elStatus.className = `${base}danger`;
      break;
    case STATUS.CLOSE:
      elStatus.className = `${base}secondary`;
      break;
    default:
      elStatus.className = `${base}warning`;
  }
};

const dataset = (label, color) => {
  return {
    label: label,
    data: [],
    backgroundColor: color,
    borderColor: color,
    pointRadius: 0,
    fill: false
  };
};

const newContainer = (name) => {
    const div = document.createElement('div');
    div.classList.add('col-6');
    div.id = name;
    elRoot.appendChild(div);

    const canvas = document.createElement('canvas');
    canvas.id = `chart#${name}`;
    div.appendChild(canvas);

    const ctx = canvas.getContext('2d');
    const config = {
        type: 'line',
        data: {
          labels: [],
          datasets: [
            dataset('mem', '#204B57'),
            dataset('cpu', '#197BBD')
          ]
        },
        options: {
          title: {
            display: true,
            text: name
          },
          scales: {
            xAxes: [{
              type: 'time',
              time: {
                displayFormats: {
                  quarter: 'h:mm:ss a'
                }
              }
            }]
          }
        }
      }
    ;

    charts.set(name, new Chart(ctx, config));
  }
;

const updateData = (data, value) => {
  if (data.length === 50) {
    data.shift();
  }
  data.push(value);
  return data;
};

const updateContainer = (name, m) => {
  const chart = charts.get(name);

  chart.data.datasets.forEach((dataset) => {
    if (dataset.label === 'mem') dataset.data = updateData(dataset.data, m.memory_percentage);
    else dataset.data = updateData(dataset.data, m.cpu_percentage);
  });
  chart.data.labels = updateData(chart.data.labels, m.time);

  chart.update();
};

const checkContainers = (metrics) => {
  metrics.forEach((m) => {
    const name = m.name;

    if (!charts.has(name)) {
      newContainer(name);
    } else {
      updateContainer(name, m);
    }
  });

  charts.forEach((val, key) => {
    let isExists = false;

    metrics.forEach((m) => {
      if (m.name === key) {
        isExists = true;
      }
    });

    if (!isExists) oldContainer(key);
  });
};

const oldContainer = (name) => {
  charts.delete(name);
  elRoot.removeChild(document.getElementById(name));
};

const start = (all = false) => {
  stop();

  socket = new WebSocket('ws://localhost:2000/metrics');

  let value = elContainersName.value;
  if (value === '' || all) value = 'all';

  socket.onopen = () => {
    setStatus(STATUS.SUCCESS);
    socket.send(value);
  };

  socket.onclose = () => {
    setStatus(STATUS.CLOSE);
  };

  socket.onmessage = (e) => {
    const data = JSON.parse(e.data);
    const alert = data.alert;
    const metrics = data.metrics;

    setAlert(alert);
    if (metrics === undefined) return;

    checkContainers(metrics);
  };

  socket.onerror = () => {
    setStatus(STATUS.ERROR);
  };
};

const stop = () => {
  setStatus(STATUS.CLOSE);

  if (socket !== null) socket.close();
  socket = null;

  charts.forEach((val, key) => {
    oldContainer(key);
  });
};
