let isModalOpenFilter = false;  // Флаг, показывающий, открыто ли модальное окно

document.querySelectorAll('.filterButton').forEach(button => {
    button.addEventListener('click', async () => {
        const table = button.closest('table');
        const columns = getTableColumns(table);  // Получаем список колонок таблицы
        const tableName = table.getAttribute('id');
        const filter = await promptSelectFilter(columns, button);  // Ждем, пока пользователь выберет фильтр
        const column = button.getAttribute('id');
        const companyID = table.getAttribute('idCompany');
        const currentProjectName = table.getAttribute('currentProjectName')
        if (!filter) {
            alert('Фильтр не выбран!');
            return;
        }

        try {
            let response;
            if (tableName === "company") {
                response = await fetch(`/filter?table=${tableName}&param=${filter.column}&column=${column}`, {
                    method: 'GET',
                    headers: { 'Content-Type': 'application/json' },
                });
            }else if(tableName === "userInCompany"){
                response = await fetch(`/filter/${companyID}?table=${tableName}&param=${filter.column}&column=${column}&companyID=${companyID}`, {
                    method: 'GET',
                    headers: { 'Content-Type': 'application/json' },
                });
            }else if(tableName === "ProjectsTask"){
                response = await fetch(`/filter/${companyID}?table=${tableName}&param=${filter.column}&column=${column}&companyID=${companyID}&projectName=${currentProjectName}`, {
                    method: 'GET',
                    headers: { 'Content-Type': 'application/json' },
                });
             } else {
                response = await fetch(`/filter/${companyID}?table=${tableName}&param=${filter.column}&column=${column}&companyID=${companyID}`, {
                    method: 'GET',
                    headers: { 'Content-Type': 'application/json' },
                });
            }

            // Получаем обновленные данные с сервера
            const data = await response.json();

            // Обновляем таблицу с новыми данными
            updateTable(data, button.closest('table')); // Передаем таблицу для обновления
        } catch (error) {
            console.error(error);
            alert('Ошибка: ' + error.message);
        }
    });
});

function getTableColumns(table) {
    const headers = table.querySelectorAll('th');
    const columns = [];
    headers.forEach(header => {
        const columnName = header.textContent.trim();
        if (columnName) {
            columns.push(columnName);
        }
    });
    return columns;
}

function promptSelectFilter(columns, button) {
    if (isModalOpenFilter) {
        alert('Окно фильтрации уже открыто');
        return Promise.resolve(null);  // Возвращаем пустое значение, чтобы не запускать процесс фильтрации
    }

    isModalOpenFilter = true;

    // Извлекаем список колонок из атрибута data-columns
    const columnList = button.getAttribute('data-columns');
    const columnOptions = columnList ? columnList.split(',') : columns;

    const modal = document.createElement('div');
    modal.classList.add('filter-modal');
    modal.style.position = 'absolute';
    modal.style.padding = '20px';
    modal.style.backgroundColor = 'white';
    modal.style.boxShadow = '0 0 10px rgba(0, 0, 0, 0.1)';
    modal.style.zIndex = '9999';
    modal.style.borderRadius = '8px';

    // Создаем селект с фильтрами для колонок из атрибута
    const selectColumn = document.createElement('select');
    columnOptions.forEach(column => {
        const option = document.createElement('option');
        option.value = column;
        option.textContent = column;
        selectColumn.appendChild(option);
    });

    const confirmButton = document.createElement('button');
    confirmButton.textContent = 'Применить фильтр';

    const closeButton = document.createElement('button');
    closeButton.textContent = 'Закрыть';

    modal.appendChild(selectColumn);
    modal.appendChild(confirmButton);
    modal.appendChild(closeButton);

    // Вставляем модальное окно в body
    document.body.appendChild(modal);

    // Получаем координаты кнопки для позиционирования окна
    const rect = button.getBoundingClientRect();
    modal.style.left = `${rect.left}px`;
    modal.style.top = `${rect.bottom + window.scrollY + 10}px`;

    return new Promise((resolve) => {
        confirmButton.addEventListener('click', () => {
            const column = selectColumn.value;
            resolve({ column });
            modal.remove();
            isModalOpenFilter = false;
        });

        closeButton.addEventListener('click', () => {
            modal.remove();
            isModalOpenFilter = false;
        });
    });
}

function updateTable(data, table) {
    const tableBody = table.querySelector('tbody');
    const tableId = table.getAttribute('id');
    const companyID = table.getAttribute('idCompany');
    const currentProjectName = table.getAttribute('currentProjectName');
    // Удаляем все строки с данными
    while (tableBody.rows.length) {
        tableBody.deleteRow(0);
    }

    // Добавляем новые строки с данными
    data.forEach((row) => {
        if (!row.hasOwnProperty('Compleate')) {
            row.Compleate = null; // Или задайте значение по умолчанию
        }
        const tr = document.createElement('tr');
        Object.entries(row).forEach(([key, value]) => {
            const td = document.createElement('td');

            switch (tableId) {
                case 'projects':
                    if (key === 'name') {
                        const link = document.createElement('a');
                        link.textContent = value;
                        // Формируем ссылку с подстановкой значений для CurrentCompanyId и Name
                        link.href = `/tasks/${companyID}/${value}`;
                        td.appendChild(link);
                    } else if (key === 'status') {
                        // Если это колонка "status", преобразуем true/false в выполненно/невыполненно
                        td.textContent = value ? 'Complete' : 'Uncompleate';
                    }else if(key ==='created'){
                        const formattedTime = formatDateTime(value);
                        td.textContent = formattedTime;
                    } else {
                        // Для остальных колонок выводим значение как есть
                        td.textContent = value;
                    }
                    break;


                    case 'ProjectsTask':
                        if (key === 'isDone') {
                            td.textContent = value ? 'Complete' : 'Uncompleate';
                        }else if(key ==='created'){
                            const formattedTime = formatDateTime(value);
                            td.textContent = formattedTime;
                        }else if(key ==='expired'){
                            const formattedTime = formatDateTime(value);
                            td.textContent = formattedTime;
                        }else if (key === 'whocomplete') {
                            td.textContent = value?.String && value.Valid ? value.String : '-';
                        }else if (key === 'Compleate') {
                            const link = document.createElement('a');
                            link.href = `/CompleteTask/${companyID}/${currentProjectName}/${row.name}`;
                            const button = document.createElement('button');
                            button.textContent = 'Complete the task';
                            link.appendChild(button);
                            td.appendChild(link);
                        }else {
                            td.textContent = value;
                        }
                        break;
            }

            tr.appendChild(td);
        });
        tableBody.appendChild(tr);
    });
}


function formatDateTime(dateString) {
    const date = new Date(dateString);

    // Форматируем дату и время с использованием метода toLocaleString
    const formattedDate = date.toLocaleString('ru-RU', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false // Используем 24-часовой формат
    });

    return formattedDate.replace(',', ''); // Убираем запятую между датой и временем
}