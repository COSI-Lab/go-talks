// Resize text area
window.onload = function () {
    let area = document.getElementById("description");
    if (area) {
        area.addEventListener("input", (e) => {
            e.target.style.height = "auto";
            e.target.style.height = (e.target.scrollHeight) + "px";
        });
    }
}

var socket = connect();
var authenticated = false;
var resolve_auth = undefined;
var reject_auth = undefined;
var password = null;

// save the password in the cookie
function setPasswordCookie(password) {
    const expirationDate = new Date();
    // expire after 30 days
    expirationDate.setDate(expirationDate.getDate() + 30);
    const cookieValue = encodeURIComponent(password);
    const cookieName = "password";
    const cookieOptions = {
        expires: expirationDate.toUTCString(),
        path: "/",
        sameSite: "strict"
    };

    document.cookie = `${cookieName}=${cookieValue}; ${Object.entries(cookieOptions)
        .map(([key, value]) => `${key}=${value}`)
        .join("; ")}`;
}

// retrieve the password from the cookie
function getPasswordCookie() {
    const cookieName = "password=";
    const decodedCookie = decodeURIComponent(document.cookie);
    const cookieArray = decodedCookie.split(';');

    // find the value with key: "password"
    for (let i = 0; i < cookieArray.length; i++) {
        let cookie = cookieArray[i].trim();
        if (cookie.startsWith(cookieName)) {
            return cookie.substring(cookieName.length);
        }
    }
    return null;
}

function removePasswordCookie() {
    // set the expiration time to a past date and then it will be deleted automatically
    authenticated = false;
    document.cookie = "password=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
}


// If the user is not authenticated trigger the auth flow
// Returns a promise that is resolved when the user is authenticated
function auth() {
    if (authenticated) {
        return Promise.resolve();
    }

    // first try to get the password from cookie, if it doesn't exist, get it from client
    if ((password = getPasswordCookie()) === null) {
        password = prompt("Please enter the password to authenticate");
    }

    const data = {
        type: 3,
        auth: {
            password: password
        }
    }

    const json = JSON.stringify(data);
    console.log(`Sending: ${json}`);

    socket.send(json);

    // Create a custom promise that is resolved by the socket onmessage handler
    const promise = new Promise((res, rej) => {
        resolve_auth = res;
        reject_auth = rej;

        // Reject after 5 seconds
        setTimeout(() => {
            rej("Timed out");
        }, 5000);
    });

    return promise;
}

// Saves the id we saw since the last sync message
var seen = new Set();
// Sync requests are used to update the table with the latest data
function sync() {
    console.log("Syncing week " + week);

    const data = {
        type: 4,
        sync: {
            week: week
        }
    }

    const json = JSON.stringify(data);
    console.log(`Sending: ${json}`);

    seen.clear();
    socket.send(json);
}

// Connects to the websocket endpoint
function connect() {
    var ws_scheme = window.location.protocol == "https:" ? "wss://" : "ws://";

    let socket = new WebSocket(ws_scheme + window.location.host + "/ws");
    socket.onopen = function (e) {
        console.log("Connected!", e);
        sync();
    };
    socket.onmessage = function (e) {
        const data = JSON.parse(e.data);
        console.log("Received:", data);

        if (data.type == 0) {
            // Add the talk to the table
            addTalk(data.new);
            seen.add(data.new.id);
        } else if (data.type == 1 || data.type == 2) {
            // Hide the talk from the table
            hideTalk(data.hide.id);
        } else if (data.type == 3) {
            // Receiving an auth message means we have successfully authenticated
            handleAuth(data.auth);
        } else if (data.type == 4) {
            // Remove talks that we didn't see in this sync
            hideTalksNotIn(seen);
            console.log("Sync complete");
        }
    };
    socket.onclose = function (e) {
        console.log("Disconnected!", e);
        socket = connect();
    };
    socket.onerror = function (e) {
        console.log("Error!", e);
    };

    return socket;
}

// User triggers talk deletion
function del(id) {
    confirmed = confirm("Are you sure you want to delete this talk?");
    if (!confirmed) {
        return;
    }

    auth().then(() => {
        const data = {
            type: 2,
            delete: {
                id: id
            }
        }

        const json = JSON.stringify(data);
        console.log(`Sending: ${json}`);

        socket.send(json);
    }).catch((reason) => {
        alert(reason);
    });
}

// User triggers talk hiding
function hide(id) {
    auth().then(() => {
        const data = {
            type: 1,
            hide: {
                id: id
            }
        };

        const json = JSON.stringify(data);
        console.log(`Sending: ${json}`);

        socket.send(json);
    }).catch((reason) => {
        alert(reason);
    });
}

// User triggers talk creation
function create() {
    const name = document.getElementById("name").value;
    const description = document.getElementById("description").value;
    const talktype = document.getElementById("type").value;

    if (name == "" || description == "" || talktype == "") {
        alert("Please fill in all fields");
        return;
    }

    auth().then(() => {
        const data = {
            type: 0,
            new: {
                name: name,
                description: description,
                talktype: typeToString[talktype],
                week: week
            }
        };

        const json = JSON.stringify(data);
        console.log(`Sending: ${json}`);

        socket.send(json);
    }).catch((reason) => {
        alert(reason);
    });
}

// ---- RESPONSE HANDLERS ---- //

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

// Removes all talks that match a predicate
function hideTalks(predicate) {
    const table = document.getElementById("tb");
    const rows = document.getElementById('tb').children;

    for (let i = rows.length - 1; i >= 0; i--) {
        if (rows[i].getAttribute("class") == "event" && predicate(rows[i])) {
            table.removeChild(rows[i]);
        }
    }
}

function hideTalk(id) {
    console.log("Hiding talk", id);

    hideTalks((row) => {
        return parseInt(row.children[0].innerText) == id;
    });
}

function hideTalksNotIn(ids) {
    console.log("Hiding talks not in", ids);

    hideTalks((row) => {
        return !ids.has(parseInt(row.children[0].innerText));
    });
}

// Removes all talks from the table
function clearTalks() {
    hideTalks((row) => {
        return true;
    });
}

function addTalk(talk) {
    if (talk.week != week) {
        console.log("Skipping new talk because it is for a different week", talk.week);
        return
    }

    const table = document.getElementById("tb");
    const rows = document.getElementById('tb').children;

    // Check to see if the talk is already in the table
    for (let i = 0; i < rows.length; i++) {
        if (rows[i].children[0].innerText == talk.id) {
            console.log("Skipping new talk because it is already in the table", talk.id);
            return
        }
    }

    // Insert the new data into the correct location in the table
    let i = 0
    for (i = 0; i < rows.length - 1; i++) {
        // Order by talk type then by id
        let childtype = stringToType[rows[i].children[2].innerText]

        if (stringToType[talk.talktype] < childtype) {
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
    c2.innerHTML = talk.talktype;

    var c3 = row.insertCell(3);
    c3.setAttribute("class", "description markdown");
    c3.innerHTML = talk.description;

    var c4 = row.insertCell(4);
    c4.setAttribute("class", "actions");
    c4.innerHTML = '<button onclick="hide(' + talk.id + ')"> x </button>';
}

// {"type": 3, "status": boolean}
function handleAuth(data) {
    success = !!data.status;

    console.log("Authenticated:", success);
    authenticated = success;

    if (resolve_auth && success) {
        setPasswordCookie(password);
        resolve_auth();
    } else if (reject_auth && !success) {
        removeAuthenticationCookie()
        reject_auth("Authentication failed");
    }
}