var editor = new FroalaEditor('#edit',{
    imageUploadURL:'/imgUpload',
    placeholderText:'在此输入内容。可以插入图片，但请不要上传文件，暂不支持!!!',
    imageManagerDeleteURL:'/imgDelete',
    language:'zh_cn',
    events: {
        'image.removed': function ($img) {
            var xhttp = new XMLHttpRequest();
            xhttp.onreadystatechange = function() {

                // Image was removed.
                if (this.readyState == 4 && this.status == 200) {
                    console.log ('image was deleted');
                }
            };
            xhttp.open("POST", "/imgDelete", true);
            xhttp.send(JSON.stringify({
                src: $img.attr('src')
            }));
        }
    }
});