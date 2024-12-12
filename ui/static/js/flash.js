function showMessage(message, messageType) {
    const messageBox = document.createElement('div');
    messageBox.textContent = message;
    messageBox.style.padding = '10px';
    messageBox.style.margin = '10px 0';
    messageBox.style.borderRadius = '5px';
    messageBox.style.color = '#fff';
    messageBox.style.fontSize = '16px';

    // Устанавливаем стиль в зависимости от типа сообщения
    switch (messageType) {
        case 'info':
            console.log(`[INFO]: ${message}`);
            messageBox.style.backgroundColor = '#007bff';
            break;
        case 'warning':
            console.log(`[WARNING]: ${message}`);
            messageBox.style.backgroundColor = '#ffc107';
            break;
        case 'error':
            console.log(`[ERROR]: ${message}`);
            messageBox.style.backgroundColor = '#dc3545';
            break;
        default:
            console.log(`[UNKNOWN]: ${message}`);
            messageBox.style.backgroundColor = '#6c757d';
            break;
    }

    // Добавляем сообщение на страницу
    document.body.appendChild(messageBox);

    // Убираем сообщение через 5 секунд
    setTimeout(() => {
        messageBox.remove();
    }, 5000);
}