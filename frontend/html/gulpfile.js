const { src, dest, watch, series, parallel } = require("gulp");
const sass = require("gulp-sass")(require("sass"));
const cleanCSS = require("gulp-clean-css");
const terser = require("gulp-terser");
const rename = require("gulp-rename");
const sourcemaps = require("gulp-sourcemaps");
const purgecss = require("gulp-purgecss");
const browserSync = require("browser-sync").create();
const fileinclude = require("gulp-file-include"); // NOVO: plugin para includes HTML

const paths = {
  html: {
    src: "./src/*.html",
    dest: "./",
  },
  styles: {
    src: "./src/scss/**/*.scss",
    dest: "./assets/css",
  },
  scripts: {
    src: "./src/js/**/*.js",
    dest: "./assets/js",
  },
  vendorFonts: {
    src: "./fonts/la-*.*",
    dest: "./assets/fonts",
  },
  purgeContent: ["./*.html", "./src/js/**/*.js", "./assets/js/**/*.js"],
};

function vendorFonts() {
  return src(paths.vendorFonts.src).pipe(dest(paths.vendorFonts.dest));
}

// NOVO: Tarefa para processar HTML com includes
function htmlDev() {
  return src(paths.html.src)
    .pipe(fileinclude({
      prefix: '@@',
      basepath: '@file',
      indent: true,  // Mantém a indentação do código
    }))
    .pipe(dest(paths.html.dest))
    .pipe(browserSync.stream());
}

// NOVO: Tarefa para build de HTML (sem sourcemaps, otimizado)
function htmlBuild() {
  return src(paths.html.src)
    .pipe(fileinclude({
      prefix: '@@',
      basepath: '@file',
      indent: false, // Remove indentação desnecessária no build
    }))
    // Opcional: minificar HTML no build
    // .pipe(htmlmin({ collapseWhitespace: true, removeComments: true }))
    .pipe(dest(paths.html.dest));
}

function stylesDev() {
  return src(paths.styles.src, { sourcemaps: true })
    .pipe(
      sass.sync({
        outputStyle: "expanded",
      }).on("error", sass.logError)
    )
    .pipe(cleanCSS())
    .pipe(
      rename({
        suffix: ".min",
      })
    )
    .pipe(dest(paths.styles.dest, { sourcemaps: "." }))
    .pipe(browserSync.stream());
}

function stylesBuild() {
  return src(paths.styles.src)
    .pipe(
      sass.sync({
        outputStyle: "expanded",
      }).on("error", sass.logError)
    )
    .pipe(
      purgecss({
        content: paths.purgeContent,
        safelist: ["is-active", "active", "open", "show"],
      })
    )
    .pipe(cleanCSS())
    .pipe(
      rename({
        suffix: ".min",
      })
    )
    .pipe(dest(paths.styles.dest));
}

function scriptsDev() {
  return src(paths.scripts.src, { sourcemaps: true })
    .pipe(
      terser({
        compress: {
          passes: 2,
        },
        mangle: false,
        format: {
          comments: false,
        },
      })
    )
    .pipe(
      rename({
        suffix: ".min",
      })
    )
    .pipe(dest(paths.scripts.dest, { sourcemaps: "." }))
    .pipe(browserSync.stream());
}

function scriptsBuild() {
  return src(paths.scripts.src)
    .pipe(
      rename({
        suffix: ".min",
      })
    )
    .pipe(
      terser({
        compress: {
          passes: 3,
          drop_console: true,
          drop_debugger: true,
          toplevel: true,
        },
        mangle: {
          toplevel: true,
        },
        format: {
          comments: false,
        },
      })
    )
    .pipe(dest(paths.scripts.dest));
}

const scriptsWatch = series(scriptsDev, reload);

function reload(done) {
  browserSync.reload();
  done();
}

function serve(done) {
  browserSync.init({
    server: {
      baseDir: "./",
    },
    notify: false,
    open: false,
  });

  done();
}

function watcher() {
  watch(paths.styles.src, stylesDev);
  watch(paths.scripts.src, scriptsWatch);
  watch(paths.vendorFonts.src, vendorFonts);
  watch(paths.html.src, htmlDev);      // ALTERADO: agora observa os HTML em src/
}

// ATUALIZADO: inclui htmlDev no desenvolvimento
const build = parallel(stylesBuild, scriptsBuild, vendorFonts, htmlBuild);
const dev = series(parallel(stylesDev, scriptsDev, vendorFonts, htmlDev), serve, watcher);

exports.styles = stylesDev;
exports.stylesBuild = stylesBuild;
exports.scriptsDev = scriptsDev;
exports.scriptsBuild = scriptsBuild;
exports.scripts = scriptsDev;
exports.vendorFonts = vendorFonts;
exports.html = htmlDev;        // NOVO: exporta tarefa HTML
exports.htmlBuild = htmlBuild; // NOVO: exporta build HTML
exports.build = build;
exports.serve = dev;
exports.default = dev;