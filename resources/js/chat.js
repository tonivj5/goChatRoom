(function() {
"use strict";
// Añadir evento al boton
var btnEnviarMsg = document.getElementById("enviarMsg");
btnEnviarMsg.addEventListener("click", enviarMsg);
// Comprobar caché del navegador
var sesion = {};

sesion.apodo = localStorage.getItem("sesion/apodo");
sesion.id = parseInt(localStorage.getItem("sesion/id"));
if(!sesion.id || !sesion.apodo) {
  sesion = {};
  sesion.apodo = prompt("Introduce tu apodo");
  sesion.id = -1;
  localStorage.setItem("sesion/apodo", sesion.apodo);
  localStorage.setItem("sesion/id", "-1");
}

// Recibir el apodo
console.log("Bienvenido al chat de Toni, "+sesion.apodo);
var inputApodo = document.getElementById('apodo');
inputApodo.innerHTML = sesion.apodo+": ";

// Obtener el ul chat
var msgs = document.getElementById('messages');

// Abrir websocket
var ws = new WebSocket("ws://localhost:5000/chat");
ws.onopen = function() {
  ws.send(JSON.stringify(sesion));
};

ws.onclose = function(e){
  console.log("Desconectado - estatus " + this.readyState);
  // TODO: Crear proceso de reconexión
//      console.log("Intentando reconectar...");
};

ws.onmessage = getMsg;

function postMsg(msg) {
  var data = {};
  var fechaActual = new Date();
  data.msg = msg.value;
  data.fecha = (fechaActual.getHours() < 10 ? "0"+fechaActual.getHours() : fechaActual.getHours()) +
                ":" + (fechaActual.getMinutes() < 10 ? "0"+fechaActual.getMinutes() : fechaActual.getMinutes());
  data = JSON.stringify(data);
  ws.send(data);
  msg.value = "";
}

function getMsg(e) {
  console.log("Mensaje recibido: ");
  console.log(e);
  var data = JSON.parse(e.data);

  switch (data.tipo) {
    case "id":
      sesion.id = data.id;
      console.log("Mi ID es: "+sesion.id);
      localStorage.setItem("sesion/id", sesion.id);

      break;
    case "msg":
      var liMsg = document.createElement("li");
      var pre = document.createElement("pre");
      pre.setAttribute("class", "msg");
      if(sesion.id == data.cliente.id)
        pre.setAttribute("class", "mio");
      pre.innerHTML = "*" + data.fecha + "* " + data.cliente.apodo + ": " + data.msg;
      liMsg.appendChild(pre);
      msgs.appendChild(liMsg);
      msgs.lastElementChild.scrollIntoView(true);

      break;
    case "conectado":
     // TODO: Debo optimizar el tratamiento del el DOM
      var conectados = [];
      for(var i = 0; i < data.clientes.length; i++) {
        var cliente = data.clientes[i];
        var liConectado = document.createElement("li");
        liConectado.setAttribute("id", cliente.id);
        liConectado.innerHTML = cliente.apodo;

        conectados[i] = liConectado;
      }

      var ulConectados = document.getElementById("users");

      for(var i = 0; i < conectados.length; i++) {
          ulConectados.appendChild(conectados[i]);
      }

      break;
    case "desconectado":
      for(var i = 0; i < data.clientes.length; i++) {
        var cliente = data.clientes[i];
        console.log(cliente);
        var liToDelete = document.getElementById(cliente.id);
        //liToDelete.parentNode.removeChild(liToDelete);
        liToDelete.setAttribute("style", "color: red;");
      }

      break;
    default:
      console.log("Tipo no identificado:");
      console.log(data);
  }
}

function enviarMsg() {
  var msg = document.getElementById("msg");
  console.log(sesion.apodo + " envía: " + msg.value);
  postMsg(msg);

  return false;
}
})();
