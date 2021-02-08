function swipedetect(el, callback) {
    var touchsurface = el,
        swipedir,
        startX,
        startY,
        distX,
        distY,
        threshold = 50, //required min distance traveled to be considered swipe
        restraint = 100, // maximum distance allowed at the same time in perpendicular direction
        allowedTime = 300, // maximum time allowed to travel that distance
        elapsedTime,
        startTime,
        handleswipe = callback || function (swipedir) {};

    touchsurface.addEventListener(
        'touchstart',
        function (e) {
            var touchobj = e.changedTouches[0];
            swipedir = 'none';
            dist = 0;
            startX = touchobj.pageX;
            startY = touchobj.pageY;
            startTime = new Date().getTime(); // record time when finger first makes contact with surface
            e.preventDefault();
        },
        false
    );

    touchsurface.addEventListener(
        'touchmove',
        function (e) {
            e.preventDefault(); // prevent scrolling when inside DIV
        },
        false
    );

    touchsurface.addEventListener(
        'touchend',
        function (e) {
            var touchobj = e.changedTouches[0];
            distX = touchobj.pageX - startX; // get horizontal dist traveled by finger while in contact with surface
            distY = touchobj.pageY - startY; // get vertical dist traveled by finger while in contact with surface
            elapsedTime = new Date().getTime() - startTime; // get time elapsed
            if (elapsedTime <= allowedTime) {
                // first condition for swipe met
                if (
                    Math.abs(distX) >= threshold &&
                    Math.abs(distY) <= restraint
                ) {
                    // 2nd condition for horizontal swipe met
                    swipedir = distX < 0 ? 'left' : 'right'; // if dist traveled is negative, it indicates left swipe
                } else if (
                    Math.abs(distY) >= threshold &&
                    Math.abs(distX) <= restraint
                ) {
                    // 2nd condition for vertical swipe met
                    swipedir = distY < 0 ? 'up' : 'down'; // if dist traveled is negative, it indicates up swipe
                } else {
                    swipedir = 'click';
                }
            }
            handleswipe(swipedir);
            e.preventDefault();
        },
        false
    );
}

function ready(callbackFunction) {
    if (document.readyState != 'loading') callbackFunction(event);
    else document.addEventListener('DOMContentLoaded', callbackFunction);
}

var darkened = {};

function darken(element) {
    element.style.opacity = 0.5;
    if (darkened[element]) clearTimeout(darkened[element]);
    darkened[element] = setTimeout(() => {
        element.style.opacity = 1.0;
    }, 300);
}

function raiseMessage(color, text) {
    var element = document.getElementById('status');
    element.innerText = text;
    switch (color) {
        case 'red':
            element.style.backgroundColor = 'rgb(80, 0, 0)';
            break;
        case 'green':
            element.style.backgroundColor = 'rgb(0, 80, 0)';
            break;
        case 'black':
            element.style.backgroundColor = 'black';
            break;
    }
}

function getVolumeAfterDelay() {
    setTimeout(() => {
        axios.post('./api/receiver/volume').then(function (resp) {
            document.getElementById('volume-label').innerText =
                resp.data.volume;
        });
    }, 1000);
}

function turnTheaterOn() {
    raiseMessage('black', 'Attempting to turn on the theater...');
    axios
        .post('./api/on')
        .then(function () {
            raiseMessage('green', 'The theater was turned on successfully.');
            getVolumeAfterDelay();
        })
        .catch(function (error) {
            raiseMessage('red', error);
            return;
        });
}

function turnTheaterOff() {
    raiseMessage('black', 'Attempting to turn off the theater...');
    axios
        .post('./api/off')
        .then(function () {
            raiseMessage('black', 'The theater was turned off successfully.');
            document.getElementById('volume-label').innerText = 'OFF';
        })
        .catch(function (error) {
            raiseMessage('red', error);
        });
}

function turnLightsOn() {
    raiseMessage('black', 'Attempting to turn on the lights...');
    axios
        .post('./api/lights/on')
        .then(function () {
            raiseMessage('green', 'The lighting was turned up to full.');
        })
        .catch(function (error) {
            raiseMessage('red', error);
        });
}

function turnLightsOff() {
    raiseMessage('black', 'Attempting to turn off the lights...');
    axios
        .post('./api/lights/off')
        .then(function () {
            raiseMessage('green', 'The lighting was turned off.');
        })
        .catch(function (error) {
            raiseMessage('red', error);
        });
}

function turnLightsToDining() {
    raiseMessage('black', 'Attempting to set the lighting to dining mode...');
    axios
        .post('./api/lights/dining')
        .then(function () {
            raiseMessage('green', 'The lighting was to dining mode.');
        })
        .catch(function (error) {
            raiseMessage('red', error);
        });
}

function turnProjectorOn() {
    raiseMessage('black', 'Attempting to turn on the projector...');
    axios
        .post('./api/projector/on')
        .then(function () {
            raiseMessage('green', 'The projector was turned on successfully.');
        })
        .catch(function (error) {
            raiseMessage('red', error);
        });
}

function turnProjectorOff() {
    raiseMessage('black', 'Attempting to turn off the projector...');
    axios
        .post('./api/projector/off')
        .then(function () {
            raiseMessage('green', 'The projector was turned off successfully.');
        })
        .catch(function (error) {
            raiseMessage('red', error);
        });
}

function turnReceiverOn() {
    raiseMessage('black', 'Attempting to turn on the receiver...');
    axios
        .post('./api/receiver/on')
        .then(function () {
            raiseMessage('green', 'The receiver was turned on successfully.');
            getVolumeAfterDelay();
        })
        .catch(function (error) {
            raiseMessage('red', error);
        });
}

function turnReceiverOff() {
    raiseMessage('black', 'Attempting to turn off the receiver...');
    axios
        .post('./api/receiver/off')
        .then(function () {
            raiseMessage('green', 'The receiver was turned off successfully.');
            document.getElementById('volume-label').innerText = 'OFF';
        })
        .catch(function (error) {
            raiseMessage('red', error);
        });
}

function turnVolumeUp() {
    axios
        .post('./api/receiver/up')
        .then(function (resp) {
            raiseMessage(
                'green',
                `The volume was increased to ${resp.data.volume}.`
            );
            document.getElementById('volume-label').innerText =
                resp.data.volume;
        })
        .catch(function (error) {
            raiseMessage('red', error);
        });
}

function turnVolumeDown() {
    axios
        .post('./api/receiver/down')
        .then(function (resp) {
            raiseMessage(
                'green',
                `The volume was decreased to ${resp.data.volume}.`
            );
            document.getElementById('volume-label').innerText =
                resp.data.volume;
        })
        .catch(function (error) {
            raiseMessage('red', error);
        });
}

function muteVolume() {
    raiseMessage('black', 'Attempting to mute the receiver...');
    axios
        .post('./api/receiver/mute')
        .then(function () {
            raiseMessage('green', 'The receiver was muted.');
            document.getElementById('volume-label').innerText = 'MUTE';
        })
        .catch(function (error) {
            raiseMessage('red', error);
        });
}

function setReceiverInput(input) {
    raiseMessage('black', `Attempting to change the input to "${input}"...`);
    axios
        .post(`./api/receiver/input/${input}`)
        .then(function () {
            raiseMessage('green', `Successfully changed input to "${input}".`);
        })
        .catch(function (error) {
            raiseMessage('red', error);
        });
}

function pressKeyOnRoku(key) {
    axios.post('./api/roku/press/' + key).catch(function () {
        raiseMessage('red', error);
    });
}

function launchOnRoku(name, appId) {
    raiseMessage('black', `Attempting to launch channel ${name}...`);
    axios
        .post('./api/roku/launch/' + appId)
        .then(function () {
            raiseMessage('green', `Successfully launched channel ${name}.`);
        })
        .catch(function (error) {
            raiseMessage('red', error);
        });
}

function sendTextToRoku(msg) {
    raiseMessage('black', `Attempting to send "${msg}" to Roku...`);
    axios
        .post('./api/roku/text', msg)
        .then(function () {
            raiseMessage('green', `Sent "${msg}" to Roku.`);
        })
        .catch(function (error) {
            raiseMessage('red', error);
        });
}

function clearTextOnRoku() {
    raiseMessage(
        'black',
        `Attempting to send 20 backspace characters to Roku...`
    );
    document.getElementById('text-input').value = '';
    axios
        .post('./api/roku/clear')
        .then(function () {
            raiseMessage('green', `Sent 20 backspace characters to Roku.`);
        })
        .catch(function (error) {
            raiseMessage('red', error);
        });
}

function submitText() {
    var msg = document.getElementById('text-input').value;
    if (msg) {
        if (
            msg.substring(0, 7) == '!input:' ||
            msg.substring(0, 7) == '$input:'
        ) {
            setReceiverInput(msg.substring(7));
        } else {
            sendTextToRoku(msg);
        }
    }
}

ready((event) => {
    // wire up the touchpad
    var touchpad = document.getElementById('touchpad');
    swipedetect(touchpad, function (swipedir) {
        if (swipedir == 'left') pressKeyOnRoku('Left');
        if (swipedir == 'right') pressKeyOnRoku('Right');
        if (swipedir == 'up') pressKeyOnRoku('Up');
        if (swipedir == 'down') pressKeyOnRoku('Down');
        if (swipedir == 'click') pressKeyOnRoku('Select');
    });
});
