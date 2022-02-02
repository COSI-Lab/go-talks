// Connects to the websocket endpoint
function connect() {
    var ws_scheme = window.location.protocol == "https:" ? "wss://" : "ws://";

    let socket = new WebSocket(ws_scheme + window.location.host + "/ws");
    socket.onopen = function(e) {
        console.log("Connected!", e);
    };
    socket.onmessage = function(e) {
        const data = JSON.parse(e.data);
        console.log("Received:", data);

        if (data.type == 0) {
            // Add the talk to the table
            addTalk(data);
        } else if (data.type == 1) {
            // Hide the talk from the table
            hideTalk(data.id);
        }
    };
    socket.onclose = function(e) {
        console.log("Disconnected!", e);
    };
    socket.onerror = function(e) {
        console.log("Error!", e);
    };
 
    return socket;
}

function hide(id) {
    const data = {
        type: 1,
        id: id
    };

    json = JSON.stringify(data);
    console.log(`Sending: ${json}`);

    socket.send(json);
}

function create() {
    const name = document.getElementById("name").value;
    const description = document.getElementById("description").value;
    const talktype = document.getElementById("type").value;

    if (name == "" || description == "" || talktype == "") {
        alert("Please fill in all fields");
        return;
    }

    const data = {
        type: 0,
        name: name,
        description: description,
        talktype: parseInt(talktype),
        week: week
    };

    json = JSON.stringify(data);
    console.log(`Sending: ${json}`);

    socket.send(json);
}

const typeToString = {
    0: "forum topic",
    1: "lightning talk",
    2: "project update",
    3: "announcement",
    4: "after-meeting slot"
}

const stringToType = {
    "forum topic": 0,
    "lightning talk": 1,
    "project update": 2,
    "announcement": 3,
    "after-meeting slot": 4
}

// Hides the talk from the table
function hideTalk(id) {
    const table = document.getElementById("tb");
    const rows = document.getElementById('tb').children;

    // Find the row to remove
    let rowToRemove = null
    for (i = 0; i < rows.length; i++) {
        if (rows[i].children[0].innerText == id) {
            rowToRemove = rows[i];
            break;
        }
    }

    if (rowToRemove) {
        // Remove the row
        table.removeChild(rowToRemove);
    }
}

// {"type":0,"name":"b","description":"b","talktype":0}
function addTalk(talk) {
    if (talk.week != week) {
        console.log("Skipping new talk because it is for a different week", talk.week);
        return
    }

    const table = document.getElementById("tb");
    const rows = document.getElementById('tb').children;

    // Insert the new data into the correct location in the table
    let i = 0
    for (i = 0; i < rows.length - 1; i++) {
        // Order by talk type then by id
        let childtype = stringToType[rows[i].children[2].innerText]

        if (talk.talktype < childtype) {
            break;
        }
    }

    // Building a new event object using _javascript_
    var row = table.insertRow(i);
    row.setAttribute("class", "event");

    var c0 = row.insertCell(0);
    c0.setAttribute("style", "display: none;");
    c0.innerHTML = talk.id;

    var c1 = row.insertCell(1);
    c1.setAttribute("class", "name");
    c1.innerHTML = talk.name;

    var c2 = row.insertCell(2);
    c2.setAttribute("class", "type");
    c2.innerHTML = typeToString[talk.talktype];

    var c3 = row.insertCell(3);
    c3.setAttribute("class", "desc");
    c3.innerHTML = talk.description;

    var c4 = row.insertCell(4);
    c4.setAttribute("class", "actions");
    c4.innerHTML = '<button onclick="hide(' + talk.id + ')"> x </button>';
}

var socket = connect();