Pod::Spec.new do |spec|
  spec.name         = 'Gely'
  spec.version      = '{{.Version}}'
  spec.license      = { :type => 'GNU Lesser General Public License, Version 3.0' }
  spec.homepage     = 'https://github.com/elysiumchain/go-elysium'
  spec.authors      = { {{range .Contributors}}
		'{{.Name}}' => '{{.Email}}',{{end}}
	}
  spec.summary      = 'iOS Elysium Client'
  spec.source       = { :git => 'https://github.com/elysiumchain/go-elysium.git', :commit => '{{.Commit}}' }

	spec.platform = :ios
  spec.ios.deployment_target  = '9.0'
	spec.ios.vendored_frameworks = 'Frameworks/Gely.framework'

	spec.prepare_command = <<-CMD
    curl https://gelystore.blob.core.windows.net/builds/{{.Archive}}.tar.gz | tar -xvz
    mkdir Frameworks
    mv {{.Archive}}/Gely.framework Frameworks
    rm -rf {{.Archive}}
  CMD
end
