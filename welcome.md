# Install punq on your cluster in your current context. This will also set up the ingress to deliver punq on your own domain. You'll be asked to confirm with "Y". 
punq install -i punq.yourdomain.com

- In your domain's DNS settings, add a record for the punq domain, e.g. punq.yourdomain.com.
- Open punq in your browser.
- Log in with the admin credentials. They are prompted to your terminal once punq is installed. Make sure to store the admin credentials in a safe place, they will only be displayed once after installation.
- The cluster where punq was installed is set up per default in your punq instance. To add more clusters, use the dropdown in the top left corner and follow the instructions. Upload your kubeconfig to add more clusters. 

**🤘 You're ready to go, have fun with punq 🤘**