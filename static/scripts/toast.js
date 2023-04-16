~function (window) {
    const notifications = document.querySelector(".notifications");
    const toastDetails = {
        timer: 3000,
        success: {
            icon: 'fa-circle-check',
        },
        error: {
            icon: 'fa-circle-xmark',
        },
        warning: {
            icon: 'fa-triangle-exclamation',
        },
        info: {
            icon: 'fa-circle-info',
        }
    }

    const removeToast = (toast) => {
        toast.classList.add("hide");
        if (toast.timeoutId) clearTimeout(toast.timeoutId);
        setTimeout(() => toast.remove(), 500);
    }
    window.removeToast = removeToast
    const createToast = (id, message) => {
        const { icon } = toastDetails[id];
        const toast = document.createElement("li");
        toast.className = `toast ${id}`;
        toast.innerHTML = `<div class="column">
                                 <i class="fa-solid ${icon}"></i>
                                 <span>${message}</span>
                              </div>
                              <i class="fa-solid fa-xmark" onclick="removeToast(this.parentElement)"></i>`;
        notifications.appendChild(toast);
        toast.timeoutId = setTimeout(() => removeToast(toast), toastDetails.timer);
    }

    window.createToast = createToast
}(window)