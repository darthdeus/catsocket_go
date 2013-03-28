set :application, "catsocket_server"

set :repository,  "git@bitbucket.org:darthdeus/catsocket_server.git"
set :scm, :git

set :user, "deploy"
set :runner, user

set :deploy_to, "/opt/apps/#{application}"
set :deploy_via, :remote_cache
set :use_sudo, false

default_run_options[:pty] = true
ssh_options[:forward_agent] = true
ssh_options[:keys] = [File.join(ENV["HOME"], ".ssh", "id_rsa")]

role :web, "catsocket.com"                          # Your HTTP server, Apache/etc
role :app, "catsocket.com"                          # This may be the same as your `Web` server

namespace :deploy do

  task :start do; end
  task :stop do; end
  task :restart do; end

end

after "deploy:create_symlink", "deploy:symlink_shared"
