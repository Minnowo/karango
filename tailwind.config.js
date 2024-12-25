module.exports = {
    darkMode: "class",
    content: ["./components/**/*.templ"],
    theme: {
        extend: {
            colors: {
                'c': {
                    "red": '#F44747',

                    "yellow": '#ffae00',
                    "orange": '#ea8500',

                    "d-black": '#181818',
                    "black"  : '#202020',

                    "d-gray": '#808080',
                    "gray"  : '#D4D4D4',

                    "white": '#eeeeee'
                },
            }
        }
    },
    corePlugins: {
        preflight: true,
    },
    experimental: {
        optimizeUniversalDefaults: true
    }
};
