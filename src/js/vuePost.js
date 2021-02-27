milu.newPost = new Vue({
    el:'#newpost',
    delimiters:['{','}'],
    data:{
        commid:0,
        err:true,
    },
    methods:{
        check:function (){
            console.log("commid",this.commid,typeof this.commid);
            if (!this.commid || this.commid <= 0) {
                this.err = true;
                return;
            }
            var res = milu.comms(this.commid);
            if (res != null)
                this.err = false;
            else
                this.err = true;
        }
    }
})