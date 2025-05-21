-- +goose Up
-- +goose StatementBegin
-- Payment Initiated
INSERT INTO templates (id, title, subject, message, type, channel, variables, metadata, created_at, updated_at)
VALUES
    -- Payment Initiated
    ('payment_initiated_email',
     'Payment Initiated Email',
     'Your payment of {{amount}} is being processed.',
     'Hi {{user_name}}, your payment of {{amount}} for {{payment_purpose}} has been initiated on {{payment_date}}. We will notify you once it is processed.',
     'transactional',
     'email',
     '[
       "user_name",
       "amount",
       "payment_purpose",
       "payment_date"
     ]',
     '{
       "language": "en"
     }',
     NOW(),
     NOW()),
-- Payment Successful
    ('payment_success_sms',
     'Payment Success SMS',
     NULL,
     'Hi {{user_name}}, your payment of {{amount}} was successful. Ref: {{payment_id}}.',
     'transactional',
     'sms',
     '[
       "user_name",
       "amount",
       "payment_id"
     ]',
     '{
       "language": "en"
     }',
     NOW(),
     NOW()),

-- Payment Failed
    ('payment_failed_email',
     'Payment Failed Email',
     'Your payment attempt failed',
     'Dear {{user_name}}, unfortunately your payment of {{amount}} failed due to {{failure_reason}}. Please try again or contact support.',
     'transactional',
     'email',
     '[
       "user_name",
       "amount",
       "failure_reason"
     ]',
     '{}',
     NOW(),
     NOW()),

-- Refund Processed
    ('refund_processed_email',
     'Refund Processed',
     'Your refund has been processed',
     'Hi {{user_name}}, we have processed your refund of {{amount}}. Refund ID: {{refund_id}}.',
     'transactional',
     'email',
     '[
       "user_name",
       "amount",
       "refund_id"
     ]',
     '{}',
     NOW(),
     NOW()),

-- OTP Sent
    ('otp_verification_sms',
     'OTP Verification SMS',
     NULL,
     'Your OTP for payment is {{otp_code}}. It will expire in 10 minutes.',
     'security',
     'sms',
     '[
       "otp_code"
     ]',
     '{}',
     NOW(),
     NOW()),

-- Payment Limit Exceeded
    ('limit_exceeded_push',
     'Limit Exceeded Push',
     NULL,
     'You have exceeded your daily payment limit of {{limit}}.',
     'warning',
     'push',
     '[
       "limit"
     ]',
     '{}',
     NOW(),
     NOW()),

-- Payment Delayed
    ('payment_delayed_email',
     'Payment Delayed Email',
     'Your payment is delayed',
     'Hi {{user_name}}, your payment of {{amount}} has been delayed. We will notify you once it is resolved.',
     'transactional',
     'email',
     '[
       "user_name",
       "amount"
     ]',
     '{}',
     NOW(),
     NOW()),

-- Subscription Renewal
    ('subscription_renewal_web',
     'Subscription Renewal Notification',
     NULL,
     'Your subscription of {{plan_name}} was renewed for {{amount}} on {{renewal_date}}.',
     'transactional',
     'web',
     '[
       "plan_name",
       "amount",
       "renewal_date"
     ]',
     '{}',
     NOW(),
     NOW()),

-- Card Expiring
    ('card_expiry_reminder_email',
     'Card Expiry Reminder',
     'Your card is expiring soon',
     'Dear {{user_name}}, your saved card ending with {{card_last4}} is expiring on {{expiry_date}}. Please update it to continue uninterrupted payments.',
     'reminder',
     'email',
     '[
       "user_name",
       "card_last4",
       "expiry_date"
     ]',
     '{}',
     NOW(),
     NOW()),

-- New Card Added
    ('card_added_sms',
     'New Card Added',
     NULL,
     'A new card ending with {{card_last4}} was added to your account.',
     'security',
     'sms',
     '[
       "card_last4"
     ]',
     '{}',
     NOW(),
     NOW()),

-- Fraud Alert
    ('fraud_alert_push',
     'Fraud Alert Push Notification',
     NULL,
     'Suspicious payment of {{amount}} detected. Please verify immediately.',
     'security',
     'push',
     '[
       "amount"
     ]',
     '{}',
     NOW(),
     NOW()),

-- Invoice Sent
    ('invoice_sent_email',
     'Invoice Sent',
     'Your invoice for {{amount}}',
     'Hi {{user_name}}, your invoice for {{amount}} is ready. Invoice ID: {{invoice_id}}.',
     'billing',
     'email',
     '[
       "user_name",
       "amount",
       "invoice_id"
     ]',
     '{}',
     NOW(),
     NOW()),

-- Wallet Balance Low
    ('wallet_low_web',
     'Low Wallet Balance Notification',
     NULL,
     'Your wallet balance is low: {{balance}}. Top up to avoid failed transactions.',
     'alert',
     'web',
     '[
       "balance"
     ]',
     '{}',
     NOW(),
     NOW()),

-- Auto Payment Scheduled
    ('auto_payment_scheduled_email',
     'Auto Payment Scheduled',
     'Auto Payment Scheduled for {{amount}}',
     'Hi {{user_name}}, your auto-payment of {{amount}} is scheduled for {{scheduled_date}}.',
     'reminder',
     'email',
     '[
       "user_name",
       "amount",
       "scheduled_date"
     ]',
     '{}',
     NOW(),
     NOW()),

-- Auto Payment Success
    ('auto_payment_success_sms',
     'Auto Payment Success SMS',
     NULL,
     'Your auto-payment of {{amount}} was successfully processed.',
     'transactional',
     'sms',
     '[
       "amount"
     ]',
     '{}',
     NOW(),
     NOW()),

-- Transaction Dispute Received
    ('dispute_received_email',
     'Transaction Dispute Received',
     'We received your dispute',
     'Hi {{user_name}}, we have received your dispute for transaction {{transaction_id}}. We will get back to you within 5 business days.',
     'support',
     'email',
     '[
       "user_name",
       "transaction_id"
     ]',
     '{}',
     NOW(),
     NOW()),

-- Dispute Resolved
    ('dispute_resolved_push',
     'Dispute Resolved Push',
     NULL,
     'Your dispute for transaction {{transaction_id}} has been resolved.',
     'support',
     'push',
     '[
       "transaction_id"
     ]',
     '{}',
     NOW(),
     NOW()),

-- Payment Confirmation Needed
    ('payment_confirmation_web',
     'Payment Confirmation Required',
     NULL,
     'Please confirm your pending payment of {{amount}} by {{confirmation_deadline}}.',
     'action_required',
     'web',
     '[
       "amount",
       "confirmation_deadline"
     ]',
     '{}',
     NOW(),
     NOW()),

-- Beneficiary Added
    ('beneficiary_added_email',
     'New Beneficiary Added',
     'You added a new beneficiary',
     'Hi {{user_name}}, you successfully added a new beneficiary: {{beneficiary_name}}.',
     'security',
     'email',
     '[
       "user_name",
       "beneficiary_name"
     ]',
     '{}',
     NOW(),
     NOW()),

-- KYC Pending
    ('kyc_pending_sms',
     'KYC Verification Pending',
     NULL,
     'Hi {{user_name}}, please complete your KYC to continue using payment services.',
     'reminder',
     'sms',
     '[
       "user_name"
     ]',
     '{}',
     NOW(),
     NOW())
ON CONFLICT (id) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE
FROM templates
WHERE id IN (
             'payment_success_sms',
             'payment_failed_email'
            'refund_processed_email',
             'otp_verification_sms',
             'limit_exceeded_push',
             'payment_delayed_email',
             'subscription_renewal_web',
             '"card_expiry_reminder_email"',
             '"card_added_sms"',
             'fraud_alert_push',
             'invoice_sent_email',
             'wallet_low_web',
             'auto_payment_scheduled_email',
             'auto_payment_success_sms'
    );
-- +goose StatementEnd
