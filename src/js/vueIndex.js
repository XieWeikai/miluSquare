var coms;
var lim = 3;
var comms = new Vue({
    el:'#comms',
    delimiters:['{','}'],
    data:{
        comunities:[],
        showAll:false,
    },
    methods:{
        showMore:function (num,clear) {
            if(clear) {
                lim = 3;
                this.showAll = false;
            }
            else
                lim += num;
            var len = lim > coms.length?coms.length:lim;
            this.comunities = coms.slice(0,len)
            if(coms.length == len)
                this.showAll = true;
        }
    },
    created:function () {
        coms = milu.comms(null,null)
        var len = lim > coms.length?coms.length:lim;
        this.comunities = coms.slice(0,len)
        if(coms.length == len)
            this.showAll = true;
    }
})

function getArti(l,r) {
    if(l >= getArti.arLen)
        return [];
    if(r > getArti.arLen - 1)
        r = getArti.arLen - 1;
    return milu.posts(l,r);
}


var posts = new Vue({
    el:'#posts',
    delimiters: ['{','}'],
    data:{
        lists:[],
        numP:10,//number of pages
        numAr:80,//number of articles
        Plist:[1,2,3],//page navigation list
        currentPage:1,
    },
    methods:{
        page:function (index) {
            //alert(index);
            this.currentPage = index;
            this.lists = getArti((index-1)*8,index*8-1);
        },
        changeP:function (st) {
            //alert(st);
            var tmp = [];
            if(st == 1){
                if(this.Plist[0]<4) {
                    //alert("1 return");
                    return;
                }
                for (key in this.Plist)
                    tmp.push(this.Plist[key] - 3);
            }else {
                if(this.Plist[0]+3>this.numP) {
                    //alert("2 return");
                    return;
                }
                for (key in this.Plist)
                    tmp.push(this.Plist[key] + 3);
            }
            this.Plist = tmp;
        }
    },
    created:function () {
        getArti.arLen = milu.posts('len');
        this.lists = getArti(0,7);
        this.numAr = getArti.arLen;
        this.numP = Math.ceil(this.numAr / 8);
    }
})