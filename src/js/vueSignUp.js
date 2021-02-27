var signUp = new Vue({
    el:'#signup',
    delimiters:['{','}'],
    data:{
        Perror:false,
        Eerror:true,
        email:'',
        password1:'',
        password2:'',
    },
    methods:{
        checkE:function () {
            var re =/^[0-9a-zA-Z\.]+@[0-9a-zA-Z]+\.[a-zA-Z]+$/;
            if (re.test(this.email))
                this.Eerror = false;
            else
                this.Eerror = true;
        },
        checkP:function () {
            if (this.password1 != this.password2)
                this.Perror = true;
            else
                this.Perror = false;
        }
    },
    computed:{
        checkPass:function () {
            if(this.password1.length < 6)
                return true;
            else
                return false;
        }
    }
})