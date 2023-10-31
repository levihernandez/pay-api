import argparse
import uuid
from sqlalchemy import create_engine, Column, String, Integer, Float, ForeignKey, DateTime
from sqlalchemy.orm import sessionmaker, declarative_base
from faker import Faker
from tqdm import tqdm
from sqlalchemy.sql import func

Base = declarative_base()


class Customer(Base):
    __tablename__ = 'customers'
    uuid = Column(String, primary_key=True)
    name = Column(String)
    last = Column(String)
    state = Column(String)
    country = Column(String)
    zip_code = Column(String)
    banking_status = Column(String)
    created_at = Column(DateTime, server_default=func.now())
    updated_at = Column(DateTime, server_default=func.now(), onupdate=func.now())


class AccountType(Base):
    __tablename__ = 'account_types'
    id = Column(Integer, primary_key=True)
    customer_uuid = Column(String, ForeignKey('customers.uuid'))
    account_type = Column(String)
    balance = Column(Float)
    created_at = Column(DateTime, server_default=func.now())
    updated_at = Column(DateTime, server_default=func.now(), onupdate=func.now())


class Ledger(Base):
    __tablename__ = 'ledger'
    id = Column(Integer, primary_key=True)
    sender_uuid = Column(String, ForeignKey('customers.uuid'))
    receiver_uuid = Column(String, ForeignKey('customers.uuid'))
    amount = Column(Float)
    account_type = Column(String)
    created_at = Column(DateTime, server_default=func.now())
    updated_at = Column(DateTime, server_default=func.now(), onupdate=func.now())


def create_and_seed_database(num_customers, num_accounts, num_payments, database_url):
    engine = create_engine(database_url)
    Base.metadata.create_all(engine)
    Session = sessionmaker(bind=engine)
    session = Session()
    faker = Faker()

    customers = []
    for _ in tqdm(range(num_customers), desc='Creating Customers'):
        customer = Customer(
            uuid=str(uuid.uuid4()),
            name=faker.first_name(),
            last=faker.last_name(),
            state=faker.state(),
            country=faker.country(),
            zip_code=faker.zipcode(),
            banking_status='Active'
        )
        customers.append(customer)
        session.add(customer)
    session.commit()

    accounts = []
    for customer in tqdm(customers, desc='Creating Accounts'):
        for _ in range(num_accounts):
            account = AccountType(
                customer_uuid=customer.uuid,
                account_type=faker.random_element(elements=('Debit', 'Credit', 'Savings')),
                balance=faker.random_number(digits=4)
            )
            accounts.append(account)
            session.add(account)
    session.commit()

    for _ in tqdm(range(num_payments), desc='Creating Payments'):
        sender = faker.random_element(elements=accounts)
        receiver = faker.random_element(elements=accounts)
        if sender != receiver and sender.balance > 0:
            amount = faker.random_number(digits=3)
            if sender.balance >= amount:
                sender.balance -= amount
                receiver.balance += amount
                ledger_entry = Ledger(
                    sender_uuid=sender.customer_uuid,
                    receiver_uuid=receiver.customer_uuid,
                    amount=amount,
                    account_type=sender.account_type
                )
                session.add(ledger_entry)
    session.commit()
    print("Database Seeding Completed!")


def main():
    parser = argparse.ArgumentParser(description='Seed a banking database.')
    parser.add_argument('--customers', type=int, default=100, help='Number of customers')
    parser.add_argument('--accounts', type=int, default=5, help='Number of accounts per customer')
    parser.add_argument('--payments', type=int, default=500, help='Number of payments')
    parser.add_argument('--db', type=str, default='sqlite:///bank.db', help='Database URL')

    args = parser.parse_args()

    create_and_seed_database(num_customers=args.customers, num_accounts=args.accounts, num_payments=args.payments,
                             database_url=args.db)


if __name__ == "__main__":
    main()
