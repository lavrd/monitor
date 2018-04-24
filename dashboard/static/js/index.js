let socket = null;

const elStatus = document.getElementById('status');
const elContainersName = document.getElementById('name');
const elAlert = document.getElementById('alert');

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

  const base = 'badge';

  switch (status) {
    case STATUS.SUCCESS:
      elStatus.classList.add(`${base}`, 'badge-success');
      break;
    case STATUS.ERROR:
      elStatus.classList.add(`${base}`, 'badge-danger');
      break;
    case STATUS.CLOSE:
      elStatus.classList.add(`${base}`, 'badge-secondary');
      break;
    default:
      elStatus.classList.add(`${base}`, 'badge-warning');
  }
};

const newContainer = () => {

};

const newContainers = (metrics) => {
  metrics.forEach((metric) => {
    console.log(metric);
  });
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

    newContainers(metrics);
  };

  socket.onerror = () => {
    setStatus(STATUS.ERROR);
  };
};

const stop = () => {
  setStatus(STATUS.CLOSE);

  if (socket !== null) socket.close();
  socket = null;
};
