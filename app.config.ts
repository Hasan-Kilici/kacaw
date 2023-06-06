export default defineAppConfig({
  docus: {
    title: 'Kacaw',
    description: 'Kacaw is a HTTP Framework for Golang',
    image: '/image-removebg-preview.png',
    socials: {
      github: 'hasan-kilici/kacaw'
    },
    github: {
      repo: 'kacaw',
      owner: 'hasan-kilici',
      edit: false
    },
    aside: {
      level: 0,
      collapsed: false,
      exclude: []
    },
    main: {
      padded: true,
      fluid: true
    },
    header: {
      logo: true,
      showLinkIcon: true,
      exclude: ['/docs']
    },
    footer: {
    }
  }
})
