-- password = "password"
INSERT INTO public.users
(id, email, "password", created_at, updated_at)
VALUES(1, 'coba@gmail.com', '$2a$10$P7dQnbh/sYrL4gAhvQble.uoxQQDYcpIK09fw50jzScwcfpptMTJG', '2025-02-22 10:49:50.475', '2025-02-22 10:49:50.475');
INSERT INTO public.users
(id, email, "password", created_at, updated_at)
VALUES(2, 'coba2@gmail.com', '$2a$10$P7dQnbh/sYrL4gAhvQble.uoxQQDYcpIK09fw50jzScwcfpptMTJG', '2025-02-22 10:49:50.475', '2025-02-22 10:49:50.475');

INSERT INTO public.wallets
(id, user_id, balance, created_at, updated_at)
VALUES(1, 1, 200000.98754309890112334300, '2025-02-22 13:13:53.886', '2025-02-22 16:24:49.314');
INSERT INTO public.wallets
(id, user_id, balance, created_at, updated_at)
VALUES(2, 2, 1000000.98787876543456787000, '2025-02-22 13:13:53.886', '2025-02-22 16:24:49.314');

INSERT INTO public.transaction_ledger
(id, wallet_id, "type", note, status, credit, debit, description, created_at)
VALUES(1, 1, 'deposit', 'test', 'completed', 200000.98754309890112334300, 0.00000000000000000000, 'test desc', '2025-02-22 15:37:04.907');
INSERT INTO public.transaction_ledger
(id, wallet_id, "type", note, status, credit, debit, description, created_at)
VALUES(2, 2, 'deposit', NULL, 'completed', 1000000.98787876543456787000, 0.00000000000000000000, NULL, '2025-02-22 16:24:49.315');