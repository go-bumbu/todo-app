function throttle(callback, delay) {
    let isThrottled = false; // Flag to control function calls
    function wrapper() {
        if (isThrottled) {
            return;
        }
        callback.apply(this, arguments);
        isThrottled = true;

        // After the delay, allow the next call
        setTimeout(function() {
            isThrottled = false;
        }, delay);
    }
    return wrapper;
}