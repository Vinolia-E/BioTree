// Show notification
export function showNotification(message, type = "info") {
    // Create notification container if it doesn't exist
    let notificationContainer = document.getElementById(
      "notification-container"
    );
    if (!notificationContainer) {
      notificationContainer = document.createElement("div");
      notificationContainer.id = "notification-container";
      document.body.appendChild(notificationContainer);
    }

    // Create notification element
    const notification = document.createElement("div");
    notification.className = `notification ${type}`;
    notification.innerHTML = `
      <span class="notification-message">${message}</span>
      <span class="notification-close">&times;</span>
    `;

    // Add to container
    notificationContainer.appendChild(notification);

    // Add close button event listener
    notification
      .querySelector(".notification-close")
      .addEventListener("click", () => {
        notification.classList.add("fade-out");
        setTimeout(() => {
          notificationContainer.removeChild(notification);
        }, 300);
      });

    // Auto-remove after 5 seconds
    setTimeout(() => {
      if (notification.parentNode === notificationContainer) {
        notification.classList.add("fade-out");
        setTimeout(() => {
          if (notification.parentNode === notificationContainer) {
            notificationContainer.removeChild(notification);
          }
        }, 300);
      }
    }, 5000);
  }