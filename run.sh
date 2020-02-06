./churndr --notification-frequency 120 \
          --noise-look-back-time 15 \
          --email-from "luislimon@gmail.com" \
          --email-to "luislimon@gmail.com" \
          --email-subject "ChurnDR: Pod with issues in monitored namespaces" \
          --smtp "smtp.gmail.com"  \
          --namespace fdp-tools-fdp-migration-service-usw2-ppd-qal \
          --namespace sandbox-sandbox-firstbank-usw2-ppd-qal \
          --namespace fdp-metadata-cgds-log-processor-usw2-ppd-e2e \
          --namespace fdp-metadata-cgds-log-processor-usw2-ppd-qal  \
          start
          
