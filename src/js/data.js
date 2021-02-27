var milu = {}

milu.INF = 2147483647;


//users(id):get user by id
//user(l,r):get user by interval
//user(l,null)  ===  user(l,INF)
//user(null,r) === user(0,r)
milu.users = function (l,r) {
    var u = null;
    var f = function (data) {
        if (data.msg === "")
            u = data.value;
    }

    if (l === 'len'){
        $.ajax('/api/user/-1', {
            dataType: 'json',
            async: false,
        }).done(f);
        return u;
    }

    if (r === undefined ) {
        $.ajax('/api/user/'+String(l),{
            dataType:'json',
            async: false
        }).done(f);
        return u;
    }

    if (l === null)
        l = 0;
    if (r === null)
        r = milu.INF;
    $.ajax('/api/user/limit',{
        dataType:'json',
        async: false,
        data:{l:l,r:r},
    }).done(f);
    return u;
}

milu.comms = function (l,r) {
    var u = null;
    var f = function(data){
        if (data.msg === "")
            u = data.value;
    }
    if(l === 'user'){
        $.ajax('/api/user/communities',{
            async:false,
            dataType:'json',
            data:{id:r}
        }).done(f)
        return u;
    }
    if (l === 'len'){
        $.ajax('/api/community/-1', {
            dataType: 'json',
            async: false,
        }).done(f);
        return u;
    }

    if(r === undefined){
        $.ajax('/api/community/'+String(l),{
            async:false,
            dataType:'json',
        }).done(f)
        return u;
    }

    if (l === null)
        l = 0;
    if (r === null)
        r = milu.INF;
    $.ajax('/api/community/limit',{
        dataType:'json',
        async: false,
        data:{l:l,r:r},
    }).done(f);
    return u;
}

milu.posts = function (l,r) {
    var u = null;
    var f = function(data){
        if (data.msg === "")
            u = data.value;
    }
    if(l === 'user'){
        $.ajax('/api/user/posts',{
            async:false,
            dataType:'json',
            data:{id:r}
        }).done(f)
        return u;
    }
    if(l === 'community'){
        $.ajax('/api/community/posts',{
            async:false,
            dataType:'json',
            data:{id:r}
        }).done(f)
        return u;
    }
    if (l === 'len'){
        $.ajax('/api/post/-1', {
            dataType: 'json',
            async: false,
        }).done(f);
        return u;
    }

    if(r === undefined){
        $.ajax('/api/post/'+String(l),{
            async:false,
            dataType:'json',
        }).done(f)
        return u;
    }


    if (l === null)
        l = 0;
    if (r === null)
        r = milu.INF;
    $.ajax('/api/post/limit',{
        dataType:'json',
        async: false,
        data:{l:l,r:r},
    }).done(f);
    return u;
}

milu.comments = function (l,r) {
    var u = null;
    var f = function(data){
        if (data.msg === "")
            u = data.value;
    }
    if(l === 'post'){
        $.ajax('/api/post/comments',{
            async:false,
            dataType:'json',
            data:{id:r}
        }).done(f)
        return u;
    }
    if (l === 'len'){
        $.ajax('/api/comment/-1', {
            dataType: 'json',
            async: false,
        }).done(f);
        return u;
    }

    if(r === undefined){
        $.ajax('/api/comment/'+String(l),{
            async:false,
            dataType:'json',
        }).done(f)
        return u;
    }

    if (l === null)
        l = 0;
    if (r === null)
        r = milu.INF;
    $.ajax('/api/comment/limit',{
        dataType:'json',
        async: false,
        data:{l:l,r:r},
    }).done(f);
    return u;
}