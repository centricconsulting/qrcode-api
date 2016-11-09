module.exports = function(grunt) {
    grunt.loadNpmTasks('grunt-git');
    grunt.loadNpmTasks('grunt-gitinfo');
    grunt.loadNpmTasks('grunt-http');
    grunt.loadNpmTasks('grunt-exec');

    grunt.initConfig({
        pkg: grunt.file.readJSON('package.json'),
        gitinfo: {
            commands: {
                'status': ['status', '--porcelain']
            }
        },
        gitlocalbranch: '<%= gitinfo.local.branch.current.name %>',
        newMaintVersion: ['<%=pkg.version.split(".")[0]%>', '<%=pkg.version.split(".")[1]%>', '<%=Number(pkg.version.split(".")[2])+1%>'].join("."),
        newMinorVersion: ['<%=pkg.version.split(".")[0]%>', '<%=Number(pkg.version.split(".")[1])+1%>', '0'].join("."),
        newMajorVersion: ['<%=Number(pkg.version.split(".")[0])+1%>', '0', '0'].join("."),
        fullVersion: '<%=pkg.version%>',

        // Add all of the changed file to the local repo.
        gitadd: {
            local: {
                options: {
                    all: true,
                    force: false
                },
                files: {
                }
            }
        }, // gitadd

        gitcommit: {
            local: {
                options: {
                    message: grunt.option('message')
                }
            }
        }, // gitcommit

        gitcheckout: {
            local: {
                options: {
                    branch: '<%= gitinfo.local.branch.current.name %>'
                }
            },
            next: {
                options: {
                    branch: "next"
                }
            },
            master: {
                options: {
                    branch: "master"
                }
            },
            test: {
                options: {
                    branch: "wjk/test"
                }
            }
        }, // git merge

        gitmerge: {
            local: {
                options: {
                    branch: '<%= gitinfo.local.branch.current.name %>',
                    squash: false
                }
            },
            next: {
                options: {
                    branch: 'next',
                    squash: false
                }
            }
        }, // git merge

        gitpush: {
            backup: {
                options: {
                    remote: "origin",
                    branch: '<%= gitinfo.local.branch.current.name %>'
                }
            },
            next: {
                options: {
                    remote: "origin",
                    branch: "next"
                }
            },
            master: {
                options: {
                    remote: "origin",
                    branch: "master",
                    tags: true
                }
            }
        }, // gitpush

        gitpull: {
            local: {
                options: {
                    remote: 'origin',
                    branch: '<%= gitinfo.local.branch.current.name %>'
                }
            },
            next: {
                options: {
                    remote: 'origin',
                    branch: 'next'
                }
            },
            master: {
                options: {
                    remote: 'origin',
                    branch: 'master'
                }
            }
        }, // gitpull

        gitrebase: {
            next: {
                options: {
                    branch: 'next'
                }
            },
            master: {
                options: {
                    branch: 'master'
                }
            }
        }, // gitrebase

        gittag: {
            local: {
                options: {
                    tag: '<%= newMaintVersion %>'
                }
            },
            next: {
                options: {
                    tag: '<%= newMaintVersion %>'
                }
            },
            demo: {
                options: {
                    tag: 'DEMO'
                }
            },
            master: {
                options: {
                    tag: '<%= fullVersion %>'
                }
            },
            maint: {
                options: {
                    tag: '<%= newMinorVersion %>'
                }
            }
        } // gittag
    }); // initConfig

    grunt.registerTask('goodmorning', ['gitinfo','gitchecklocal','gitcheckout:next','gitpull:next','gitcheckout:local','gitrebase:next','test','jshint:all']);
    grunt.registerTask('goodnight', ['gitinfo','gitchecklocal','gitpush:backup']);

    // Checkout the "next" branch.
    grunt.registerTask('gitchecklocal', 'Check that local code is ready to go...', function() {
        grunt.task.requires('gitinfo');
        var status = grunt.config('gitinfo.status');

        if (status !== '') {
          grunt.fail.fatal(status + ': There are uncommitted local modifications.');
          return false;
        }
    });

    // Build a local copy.
    grunt.registerTask('build', 'Build locally...', function() {
        if((grunt.option('message') === "") || (grunt.option('message') === undefined)) {
            grunt.log.error("REQD: --message='<git Commit Message>'");
            grunt.fail.fatal("Git commit message is required.", 3);
            return;
        } //
        grunt.config.requires('pkg');
        // Run the build.
        grunt.task.run(['gitadd:local','gitcommit:local']);
    });

    // Checkout the "next" branch.
    grunt.registerTask('build2demo', 'Building to demo environment...', function() {
        grunt.config.requires('pkg');
        // Increment the build version number.
        grunt.config.set('pkg.version', grunt.config.get('newMaintVersion'));
        // Write the config file back with the new changes.
        grunt.file.write('package.json', JSON.stringify(grunt.config.get('pkg'), null, 2));
        // Run the build.
        grunt.task.run(['build','gitinfo','gitcheckout:next','gitpull:next','gitmerge:local','gitpush:next','gitcheckout:local']);
    });

    // Build the production environment.
    grunt.registerTask('build2prod', 'Build to production environment...', function() {
        // Run the build.
        grunt.task.run(['gitinfo','gitcheckout:master','gitpull:master','gitmerge:next','gittag:master','gitpush:master','gitcheckout:local']);
    });

    // Upgrade and put in to demo.
    grunt.registerTask('upgrade', 'Upgrade the demo environment...', function() {
        // Increment the minor version number.
        grunt.config.set('pkg.version', grunt.config.get('newMinorVersion'));
        // Write the config file back with the new changes.
        grunt.file.write('package.json', JSON.stringify(grunt.config.get('pkg'), null, 2));
        // Run the build.
        grunt.task.run(['build','gitinfo','gitchecklocal','gitcheckout:next','gitpull:next','gitmerge:local','gitpush:next','gitcheckout:local']);
    });

}; // Gruntfile
