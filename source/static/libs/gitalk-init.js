var gitalk = new Gitalk({
    clientID: 'd7b7061e43e1334a5aa9',
    clientSecret: '5af2efcfc06b854bc5ccc23babfa3a06a7508c1e',
    id: md5(window.location.pathname),
    repo: 'blog',
    owner: 'wangerzi',
    admin: 'wangerzi'
})
gitalk.render('gitalk')