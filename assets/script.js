/* placeholder file for JavaScript */

const confirm_delete = (id) => {
    if (window.confirm(`Task ${id} を削除します．よろしいですか？`)) {
        location.href = `/task/delete/${id}`;
    }
}

const confirm_update = (id) => {
    if (window.confirm(`Task ${id} を編集します．よろしいですか？`)) {
        document.edit.submit();
    }
}
