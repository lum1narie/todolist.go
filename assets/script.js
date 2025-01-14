/* placeholder file for JavaScript */

const confirm_task_delete = (id) => {
    if (window.confirm(`Task ${id} を削除します。よろしいですか？`)) {
        location.href = `/task/delete/${id}`;
    }
}

const confirm_user_delete = () => {
    if (window.confirm(`ログイン中のユーザーを削除します。よろしいですか？`)) {
        fetch(`/user/me`, {
            method: 'DELETE',
        })
            .then((response) => {
                console.log(response)
                if (response.redirected) {
                    location.href = response.url;
                }
            })
    }
}

const confirm_task_update = (id) => {
    if (window.confirm(`Task ${id} を編集します。よろしいですか？`)) {
        document.edit.submit();
    }
}
const confirm_user_update = () => {
    if (window.confirm(`ログイン中のユーザーを編集します。よろしいですか？`)) {
        document.edit.submit();
    }
}
