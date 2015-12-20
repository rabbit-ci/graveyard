import {
  moduleForModel,
  test
} from 'ember-qunit';

moduleForModel('script', {
  // Specify the other units that are required for this test.
  needs: ['model:stdio']
});

test('it exists', function(assert) {
  var model = this.subject();
  // var store = this.store();
  assert.ok(!!model);
});
