module.exports = function(grunt) {
  grunt.initConfig({
    pkg: grunt.file.readJSON('package.json'),
    jshint: {
      // define the files to lint
      files: ['Gruntfile.js', 'app/*.js'],
      // configure JSHint (documented at http://www.jshint.com/docs/)
      options: {
          // more options here if you want to override JSHint defaults
        globals: {
          console: true
        }
      }
    },
    less: {
      dev: {
        options: {
          paths: ['app']
        },
        files: {
          'built/style.css': ['app/*.less']
        }
      }
    },
    watch: {
      files: ['<%= jshint.files %>', 'app/*.less'],
      tasks: ['less:dev'],
      options: {
        interrupt: true
      }
    },
    connect: {
      server: {
        options: {
          hostname: "*"
        }
      }
    }
  });

  grunt.loadNpmTasks('grunt-contrib-watch');
  grunt.loadNpmTasks('grunt-contrib-connect');
  grunt.loadNpmTasks('grunt-contrib-jshint');
  grunt.loadNpmTasks('grunt-contrib-less');

  grunt.registerTask('dev', ['connect', 'watch']);
};
