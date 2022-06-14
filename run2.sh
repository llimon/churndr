./churndr2 --notification-frequency 120 \
          --port 8081 \
          --noise-look-back-time 15 \
          --email-from "luislimon@gmail.com" \
          --email-to "luislimon@gmail.com" \
          --email-subject "ChurnDR: Pod with issues in monitored namespaces" \
          --smtp "smtp.gmail.com"  \
          --namespace default \
          --namespace dev-devx-datacollector-usw2-prd \
          start
          
