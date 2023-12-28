"""update dalle add image_url.

Revision ID: 3d7d9e74d35d
Revises: fb32cb00c20d
Create Date: 2023-12-19 14:57:39.439042

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = '3d7d9e74d35d'
down_revision = 'fb32cb00c20d'
branch_labels = None
depends_on = None


def upgrade() -> None:
    # ### commands auto generated by Alembic - please adjust! ###
    op.add_column('dalle', sa.Column('image_url', sa.String(length=200), nullable=True))
    # ### end Alembic commands ###


def downgrade() -> None:
    # ### commands auto generated by Alembic - please adjust! ###
    op.drop_column('dalle', 'image_url')
    # ### end Alembic commands ###
